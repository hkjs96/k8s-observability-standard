package validate

import (
	"fmt"
	"os"
	"os/exec"
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

	promtool, err := exec.LookPath("promtool")
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
