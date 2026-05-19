package walk

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestFilesSkipsLocalWorkDirs(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "kept.yaml")
	writeTestFile(t, root, ".cache/ignored.yaml")
	writeTestFile(t, root, ".git/ignored.yaml")
	writeTestFile(t, root, ".tmp/ignored.yaml")

	files, err := Files(root, func(path string) bool {
		return filepath.Ext(path) == ".yaml"
	})
	if err != nil {
		t.Fatal(err)
	}

	for i, file := range files {
		files[i], err = filepath.Rel(root, filepath.FromSlash(file))
		if err != nil {
			t.Fatal(err)
		}
		files[i] = filepath.ToSlash(files[i])
	}
	sort.Strings(files)

	want := []string{"kept.yaml"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("Files() = %v, want %v", files, want)
	}
}

func TestHasPrefix(t *testing.T) {
	if !HasPrefix(filepath.Join("values", "profiles", "basic.yaml"), "values/") {
		t.Fatal("HasPrefix() = false, want true")
	}
	if HasPrefix("README.md", "values/") {
		t.Fatal("HasPrefix() = true, want false")
	}
}

func writeTestFile(t *testing.T, root, name string) {
	t.Helper()

	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("test: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}
