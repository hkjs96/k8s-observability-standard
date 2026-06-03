package validate

import (
	"strings"
	"testing"
)

const testChartLock = `apiVersion: observability.standard/v1
kind: ChartVersions
metadata:
  name: test
spec:
  reviewedAt: "2026-01-01"
  charts:
    kube-prometheus-stack:
      repository: https://prometheus-community.github.io/helm-charts
      chart: kube-prometheus-stack
      version: 85.0.2
      appVersion: v0.90.1
    loki:
      version: 7.0.0
    alloy:
      version: 1.8.2
    tempo:
      version: 1.24.4
`

func argoApp(revision string) string {
	return `apiVersion: argoproj.io/v1alpha1
kind: Application
spec:
  sources:
    - repoURL: https://prometheus-community.github.io/helm-charts
      chart: kube-prometheus-stack
      targetRevision: ` + revision + `
`
}

func TestChartVersionsParsesLock(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "charts-lock/chart-versions.yaml", testChartLock)

		versions, err := chartVersions()
		if err != nil {
			t.Fatalf("chartVersions() error = %v", err)
		}

		want := map[string]string{
			"kube-prometheus-stack": "85.0.2",
			"loki":                  "7.0.0",
			"alloy":                 "1.8.2",
			"tempo":                 "1.24.4",
		}
		for chart, version := range want {
			if versions[chart] != version {
				t.Errorf("chartVersions()[%q] = %q, want %q", chart, versions[chart], version)
			}
		}
		// appVersion must not be picked up as the chart version.
		if versions["kube-prometheus-stack"] != "85.0.2" {
			t.Errorf("appVersion leaked into chart version: %q", versions["kube-prometheus-stack"])
		}
	})
}

func TestChartsConsistent(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "charts-lock/chart-versions.yaml", testChartLock)
		writeFile(t, root, "argocd/applications/10-kube-prometheus-stack.yaml", argoApp("85.0.2"))

		if err := Charts(); err != nil {
			t.Fatalf("Charts() error = %v, want nil", err)
		}
	})
}

func TestChartsDetectsArgoDrift(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "charts-lock/chart-versions.yaml", testChartLock)
		writeFile(t, root, "argocd/applications/10-kube-prometheus-stack.yaml", argoApp("99.9.9"))

		err := Charts()
		if err == nil {
			t.Fatal("Charts() error = nil, want targetRevision drift error")
		}
		if !strings.Contains(err.Error(), "targetRevision") {
			t.Fatalf("Charts() error = %v, want targetRevision detail", err)
		}
	})
}

func TestChartsRequiresEveryPinnedChart(t *testing.T) {
	withWorkDir(t, func(root string) {
		lockWithoutTempo := strings.Replace(testChartLock, "    tempo:\n      version: 1.24.4\n", "", 1)
		writeFile(t, root, "charts-lock/chart-versions.yaml", lockWithoutTempo)
		writeFile(t, root, "argocd/applications/10-kube-prometheus-stack.yaml", argoApp("85.0.2"))

		err := Charts()
		if err == nil {
			t.Fatal("Charts() error = nil, want missing chart error")
		}
		if !strings.Contains(err.Error(), "tempo") {
			t.Fatalf("Charts() error = %v, want tempo detail", err)
		}
	})
}
