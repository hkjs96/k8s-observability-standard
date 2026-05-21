package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindPromtoolUsesLocalTmpTool(t *testing.T) {
	withWorkDir(t, func(root string) {
		want := filepath.Join(".tmp", "tools", "promtool.exe")
		writeFile(t, root, want, "fake promtool")

		got, err := findPromtool()
		if err != nil {
			t.Fatalf("findPromtool() error = %v", err)
		}
		if got != want {
			t.Fatalf("findPromtool() = %q, want %q", got, want)
		}
	})
}

func TestFindPromtoolUsesExpandedPrometheusArchive(t *testing.T) {
	withWorkDir(t, func(root string) {
		want := filepath.Join(".tmp", "tools", "prometheus-3.4.0.windows-amd64", "promtool.exe")
		writeFile(t, root, want, "fake promtool")

		got, err := findPromtool()
		if err != nil {
			t.Fatalf("findPromtool() error = %v", err)
		}
		if got != want {
			t.Fatalf("findPromtool() = %q, want %q", got, want)
		}
	})
}

func TestPrometheusRulesRequiresPromtoolMirror(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "rules/alerting/basic-alerts.yaml", "apiVersion: monitoring.coreos.com/v1\n")

		err := PrometheusRules(Options{})
		if err == nil {
			t.Fatal("PrometheusRules() error = nil, want missing mirror error")
		}
		if !strings.Contains(err.Error(), "missing promtool mirror") {
			t.Fatalf("PrometheusRules() error = %v, want missing mirror detail", err)
		}
	})
}

func TestPrometheusRulesAllowsEmptyRulesDirectory(t *testing.T) {
	withWorkDir(t, func(root string) {
		if err := os.MkdirAll("rules", 0o755); err != nil {
			t.Fatal(err)
		}

		if err := PrometheusRules(Options{}); err != nil {
			t.Fatalf("PrometheusRules() error = %v", err)
		}
	})
}
