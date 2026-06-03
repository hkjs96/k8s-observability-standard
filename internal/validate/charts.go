package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const kubePrometheusStack = "kube-prometheus-stack"

// chartVersions reads charts-lock/chart-versions.yaml and returns the pinned
// version for each chart keyed by chart name. charts-lock is the single source
// of truth for chart versions; the validators and the rendered Argo CD
// Application must agree with it. The lock file is repository-owned and simple,
// so it is parsed with a small indentation scanner to keep obsctl dependency
// free (matching the rest of internal/validate).
func chartVersions() (map[string]string, error) {
	path := filepath.Join("charts-lock", "chart-versions.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read chart lock: %w", err)
	}

	versions := map[string]string{}
	inCharts := false
	current := ""
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		switch {
		case indent <= 2 && trimmed == "charts:":
			inCharts = true
			current = ""
		case indent <= 2:
			inCharts = false
		case !inCharts:
			// outside the charts block
		case indent == 4 && strings.HasSuffix(trimmed, ":"):
			current = strings.TrimSuffix(trimmed, ":")
		case indent >= 6 && current != "" && strings.HasPrefix(trimmed, "version:"):
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "version:"))
			versions[current] = strings.Trim(value, `"'`)
		}
	}
	return versions, nil
}

// requireChartVersions loads the chart lock and errors if any of the named
// charts is missing a pinned version.
func requireChartVersions(names ...string) (map[string]string, error) {
	versions, err := chartVersions()
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if versions[name] == "" {
			return nil, fmt.Errorf("charts-lock/chart-versions.yaml has no version for chart %q", name)
		}
	}
	return versions, nil
}

// lockedChartVersion returns the pinned version for a single chart, or an error
// if charts-lock does not pin it.
func lockedChartVersion(name string) (string, error) {
	versions, err := requireChartVersions(name)
	if err != nil {
		return "", err
	}
	return versions[name], nil
}

// Charts verifies that pinned chart versions stay consistent across the
// repository. charts-lock is the source of truth, and the rendered Argo CD
// Application must reference the same kube-prometheus-stack version so the
// validated render and the deployed render cannot drift apart.
func Charts() error {
	versions, err := requireChartVersions(kubePrometheusStack, "loki", "alloy", "tempo")
	if err != nil {
		return err
	}

	appPath := filepath.Join("argocd", "applications", "10-kube-prometheus-stack.yaml")
	revision, err := argoChartRevision(appPath, kubePrometheusStack)
	if err != nil {
		return err
	}
	if revision != versions[kubePrometheusStack] {
		return fmt.Errorf(
			"argocd Application targetRevision %q does not match charts-lock version %q for %s",
			revision, versions[kubePrometheusStack], kubePrometheusStack,
		)
	}

	fmt.Println("charts ok")
	return nil
}

// argoChartRevision returns the targetRevision of the source whose chart field
// matches the given chart name in an Argo CD Application manifest.
func argoChartRevision(path, chart string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	matched := false
	for _, raw := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(raw)
		if strings.HasPrefix(trimmed, "chart:") {
			matched = strings.TrimSpace(strings.TrimPrefix(trimmed, "chart:")) == chart
		}
		if matched && strings.HasPrefix(trimmed, "targetRevision:") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "targetRevision:")), nil
		}
	}
	return "", fmt.Errorf("targetRevision for chart %q not found in %s", chart, path)
}
