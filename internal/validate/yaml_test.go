package validate

import (
	"os/exec"
	"strings"
	"testing"
)

func TestYAMLRejectsInvalidYAML(t *testing.T) {
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python unavailable")
	}

	withWorkDir(t, func(root string) {
		writeFile(t, root, "values/bad.yaml", "key: [unterminated\n")

		err := YAML()
		if err == nil {
			t.Fatal("YAML() error = nil, want parse error")
		}
		if !strings.Contains(err.Error(), "YAML validation failed") {
			t.Fatalf("YAML() error = %v, want validation detail", err)
		}
	})
}

func TestYAMLAllowsValidYAML(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "values/good.yaml", "key: value\n")

		if err := YAML(); err != nil {
			t.Fatalf("YAML() error = %v", err)
		}
	})
}
