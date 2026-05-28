package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckLogLabelCardinalityAcceptsProfileValues(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, root, "values/profiles/logs.yaml", "loki: {}\n")
	mustWrite(t, root, "values/profiles/logs-alloy.yaml", "target_label = \"namespace\"\n")
	mustWrite(t, root, "examples/phase3-smoke/log-generator.yaml", "app: example-log-generator\n")
	withWorkingDir(t, root, func() {
		if err := checkLogLabelCardinality(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCheckLogLabelCardinalityRejectsHighCardinalityLabels(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, root, "values/profiles/logs.yaml", "loki: {}\n")
	mustWrite(t, root, "values/profiles/logs-alloy.yaml", "__meta_kubernetes_pod_uid\n")
	mustWrite(t, root, "examples/phase3-smoke/log-generator.yaml", "app: example-log-generator\n")
	withWorkingDir(t, root, func() {
		if err := checkLogLabelCardinality(); err == nil {
			t.Fatal("checkLogLabelCardinality() error = nil, want high-cardinality label error")
		}
	})
}

func mustWrite(t *testing.T, root, path, content string) {
	t.Helper()
	full := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func withWorkingDir(t *testing.T, dir string, fn func()) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatal(err)
		}
	}()
	fn()
}
