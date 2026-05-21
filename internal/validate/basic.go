package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Basic(opts Options) error {
	helm, err := findHelm()
	if err != nil {
		if opts.StrictTools {
			return err
		}
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
	rendered, err := exec.Command(helm, templateArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("helm template failed: %w\n%s", err, string(rendered))
	}

	kubeconform, err := findKubeconform()
	if err != nil {
		if opts.StrictTools {
			return err
		}
		fmt.Println("kubeconform unavailable; skipped rendered manifest validation")
	} else if err := runKubeconform(kubeconform, rendered); err != nil {
		return err
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

func findKubeconform() (string, error) {
	candidates := []string{
		filepath.Join(".tmp", "tools", "kubeconform.exe"),
		filepath.Join(".tmp", "tools", "kubeconform"),
		filepath.Join("..", "lgtm-k8s-observability-v2", "tools", "bin", "kubeconform.exe"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	if kubeconform, err := exec.LookPath("kubeconform"); err == nil {
		return kubeconform, nil
	}

	return "", fmt.Errorf("kubeconform not found")
}

func runKubeconform(kubeconform string, rendered []byte) error {
	cmd := exec.Command(kubeconform, "-strict", "-ignore-missing-schemas", "-summary")
	cmd.Stdin = strings.NewReader(string(rendered))
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	if err != nil {
		return fmt.Errorf("kubeconform failed: %w", err)
	}
	return nil
}
