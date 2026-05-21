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
