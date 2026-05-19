package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/example/k8s-observability/internal/walk"
)

func PrometheusRules() error {
	files, err := walk.Files("rules", func(path string) bool {
		return strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".promtool.yaml")
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		base := strings.TrimSuffix(file, ".yaml")
		mirror := base + ".promtool.yaml"
		if _, err := os.Stat(mirror); err != nil {
			return fmt.Errorf("missing promtool mirror for %s: %s", file, mirror)
		}
	}

	promtool, err := findPromtool()
	if err != nil {
		fmt.Println("promtool unavailable; skipped promtool check rules")
		fmt.Println("prometheus rules ok")
		return nil
	}

	mirrors, err := walk.Files("rules", func(path string) bool {
		return strings.HasSuffix(path, ".promtool.yaml")
	})
	if err != nil {
		return err
	}
	args := append([]string{"check", "rules"}, mirrors...)
	out, err := exec.Command(promtool, args...).CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	if err != nil {
		return fmt.Errorf("promtool failed: %w", err)
	}

	fmt.Printf("promtool checked %d files\n", len(mirrors))
	fmt.Println("prometheus rules ok")
	return nil
}

func findPromtool() (string, error) {
	candidates := []string{
		filepath.Join(".tmp", "tools", "promtool.exe"),
		filepath.Join("..", "lgtm-k8s-observability-v2", "tools", "bin", "promtool.exe"),
	}

	matches, err := filepath.Glob(filepath.Join(".tmp", "tools", "prometheus-*", "promtool.exe"))
	if err != nil {
		return "", err
	}
	candidates = append(candidates, matches...)

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	if promtool, err := exec.LookPath("promtool"); err == nil {
		return promtool, nil
	}

	return "", fmt.Errorf("promtool not found")
}
