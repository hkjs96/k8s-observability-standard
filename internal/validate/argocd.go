package validate

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/example/k8s-observability/internal/walk"
)

func ArgoCD() error {
	files, err := walk.Files("argocd", func(path string) bool {
		return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
	})
	if err != nil {
		return err
	}

	defaultProject := regexp.MustCompile(`project:\s+default\b`)
	wildcard := regexp.MustCompile(`sourceRepos:\s*\['\*'\]|namespace:\s*'\*'`)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		text := string(data)
		if defaultProject.MatchString(text) {
			return fmt.Errorf("Argo CD Application must not use default project: %s", file)
		}
		if wildcard.MatchString(text) {
			return fmt.Errorf("Argo CD wildcard sourceRepos or destination namespace found: %s", file)
		}
		if strings.Contains(text, "chart:") && !strings.Contains(text, "valueFiles:") {
			return fmt.Errorf("Helm Application missing valueFiles: %s", file)
		}
	}

	fmt.Println("argocd ok")
	return nil
}
