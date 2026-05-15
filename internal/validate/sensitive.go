package validate

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/example/k8s-observability/internal/walk"
)

func SensitiveValues() error {
	files, err := walk.Files(".", func(path string) bool {
		if path == "README.md" {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

	forbidden := regexp.MustCompile(`adminPassword|access_key|secret_key|arn:aws|sourceRepos:\s*\['\*'\]|namespace:\s*'\*'|NodePort|admin123|shi-cluster|loki-chunks|ACCESS_KEY|SECRET_KEY`)
	var hits []string
	for _, file := range files {
		if file == "internal/validate/sensitive.go" || file == "scripts/validate-sensitive-values.ps1" {
			continue
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(data), "\n") {
			if forbidden.MatchString(line) {
				hits = append(hits, fmt.Sprintf("%s:%d:%s", file, i+1, line))
			}
		}
	}
	if len(hits) > 0 {
		return fmt.Errorf("forbidden sensitive or implementation-specific pattern found:\n%s", strings.Join(hits, "\n"))
	}

	allowedCustomers := map[string]bool{
		"AGENTS.md":                             true,
		".agent/checks/sensitive-values.md":     true,
		".agent/rules/repository-boundary.md":   true,
		"docs/00-overview.md":                   true,
		"scripts/validate-sensitive-values.ps1": true,
		"internal/validate/sensitive.go":        true,
	}
	var customerHits []string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), "customers/") && !allowedCustomers[file] {
			customerHits = append(customerHits, file)
		}
	}
	if len(customerHits) > 0 {
		return fmt.Errorf("unexpected customers/ mention found: %s", strings.Join(customerHits, ", "))
	}

	fmt.Println("sensitive values ok")
	return nil
}
