package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSensitiveValuesSkipsIgnoredWorkDirs(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, ".tmp/tools/prometheus.exe", "AccessKeyID "+"SECRET"+"_KEY")
		writeFile(t, root, ".cache/build/output.txt", "admin"+"Password")
		writeFile(t, root, "values/profiles/basic.yaml", "cluster: example-cluster\n")

		if err := SensitiveValues(); err != nil {
			t.Fatalf("SensitiveValues() error = %v", err)
		}
	})
}

func TestSensitiveValuesRejectsForbiddenPattern(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "values/profiles/basic.yaml", "admin"+"Password: admin"+"123\n")

		err := SensitiveValues()
		if err == nil {
			t.Fatal("SensitiveValues() error = nil, want forbidden pattern error")
		}
		if !strings.Contains(err.Error(), "values/profiles/basic.yaml") {
			t.Fatalf("SensitiveValues() error = %v, want file path", err)
		}
	})
}

func TestSensitiveValuesRejectsCustomersOutsideAllowedDocs(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "docs/01-standard-decisions.md", "Use customers"+"/ here.\n")

		err := SensitiveValues()
		if err == nil {
			t.Fatal("SensitiveValues() error = nil, want customers path error")
		}
		if !strings.Contains(err.Error(), "docs/01-standard-decisions.md") {
			t.Fatalf("SensitiveValues() error = %v, want file path", err)
		}
	})
}

func TestSensitiveValuesRejectsArgoWildcardPattern(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "argocd/projects/project.yaml", "source"+"Repos: ['*']\n")

		err := SensitiveValues()
		if err == nil {
			t.Fatal("SensitiveValues() error = nil, want wildcard error")
		}
		if !strings.Contains(err.Error(), "source"+"Repos") {
			t.Fatalf("SensitiveValues() error = %v, want sourceRepos detail", err)
		}
	})
}

func withWorkDir(t *testing.T, fn func(root string)) {
	t.Helper()

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Fatal(err)
		}
	})

	fn(root)
}

func writeFile(t *testing.T, root, name, content string) {
	t.Helper()

	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
