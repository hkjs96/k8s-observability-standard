package validate

import (
	"path/filepath"
	"testing"
)

func TestFindKubeconformUsesLocalTmpTool(t *testing.T) {
	withWorkDir(t, func(root string) {
		want := filepath.Join(".tmp", "tools", "kubeconform.exe")
		writeFile(t, root, want, "fake kubeconform")

		got, err := findKubeconform()
		if err != nil {
			t.Fatalf("findKubeconform() error = %v", err)
		}
		if got != want {
			t.Fatalf("findKubeconform() = %q, want %q", got, want)
		}
	})
}

func TestFindHelmUsesLocalTmpTool(t *testing.T) {
	withWorkDir(t, func(root string) {
		want := filepath.Join(".tmp", "tools", "helm")
		writeFile(t, root, want, "fake helm")

		got, err := findHelm()
		if err != nil {
			t.Fatalf("findHelm() error = %v", err)
		}
		if got != want {
			t.Fatalf("findHelm() = %q, want %q", got, want)
		}
	})
}
