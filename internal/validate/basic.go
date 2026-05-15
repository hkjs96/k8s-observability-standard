package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Basic() error {
	helm, err := findHelm()
	if err != nil {
		fmt.Println("helm unavailable; skipped helm template and helm lint")
		return nil
	}

	templateArgs := []string{
		"template", "kube-prometheus-stack", "kube-prometheus-stack",
		"--repo", "https://prometheus-community.github.io/helm-charts",
		"--version", "85.0.2",
		"--namespace", "monitoring",
		"-f", "values/common/kube-prometheus-stack.yaml",
		"-f", "values/profiles/basic.yaml",
		"-f", "values/env/dev.yaml",
		"-f", "values/sizing/small.yaml",
	}
	if out, err := exec.Command(helm, templateArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("helm template failed: %w\n%s", err, string(out))
	}

	lintDir := filepath.Join(".cache", fmt.Sprintf("helm-lint-go-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(lintDir, 0o755); err != nil {
		return err
	}
	pullArgs := []string{
		"pull", "kube-prometheus-stack",
		"--repo", "https://prometheus-community.github.io/helm-charts",
		"--version", "85.0.2",
		"--untar",
		"--untardir", lintDir,
	}
	if out, err := exec.Command(helm, pullArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("helm pull failed: %w\n%s", err, string(out))
	}

	chartPath := filepath.Join(lintDir, "kube-prometheus-stack")
	lintArgs := []string{
		"lint", chartPath,
		"-f", "values/common/kube-prometheus-stack.yaml",
		"-f", "values/profiles/basic.yaml",
		"-f", "values/env/dev.yaml",
		"-f", "values/sizing/small.yaml",
	}
	out, err := exec.Command(helm, lintArgs...).CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	if err != nil {
		return fmt.Errorf("helm lint failed: %w", err)
	}

	fmt.Println("basic validation ok")
	return nil
}

func findHelm() (string, error) {
	if helm, err := exec.LookPath("helm"); err == nil {
		return helm, nil
	}
	local := filepath.Join("..", "lgtm-k8s-observability-v2", "tools", "bin", "helm.exe")
	if _, err := os.Stat(local); err == nil {
		return local, nil
	}
	return "", fmt.Errorf("helm not found")
}
