package validate

import (
	"os"
	"strings"
	"testing"
)

func TestArgoCDRejectsWildcardSourceRepos(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "argocd/projects/project.yaml", "source"+"Repos: ['*']\n")

		err := ArgoCD()
		if err == nil {
			t.Fatal("ArgoCD() error = nil, want wildcard error")
		}
		if !strings.Contains(err.Error(), "wildcard") {
			t.Fatalf("ArgoCD() error = %v, want wildcard detail", err)
		}
	})
}

func TestArgoCDRequiresHelmValueFiles(t *testing.T) {
	withWorkDir(t, func(root string) {
		writeFile(t, root, "argocd/applications/app.yaml", "project: observability\nchart: kube-prometheus-stack\n")

		err := ArgoCD()
		if err == nil {
			t.Fatal("ArgoCD() error = nil, want missing valueFiles error")
		}
		if !strings.Contains(err.Error(), "valueFiles") {
			t.Fatalf("ArgoCD() error = %v, want valueFiles detail", err)
		}
	})
}

func TestArgoCDAllowsEmptyDirectory(t *testing.T) {
	withWorkDir(t, func(root string) {
		if err := os.Mkdir("argocd", 0o755); err != nil {
			t.Fatal(err)
		}

		if err := ArgoCD(); err != nil {
			t.Fatalf("ArgoCD() error = %v", err)
		}
	})
}
