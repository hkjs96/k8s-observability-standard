package validate

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Traces(opts Options) error {
	helm, err := findHelm()
	if err != nil {
		if opts.StrictTools {
			return err
		}
		fmt.Println("helm unavailable; skipped traces helm template")
		return checkTracePlaceholders()
	}

	args := []string{
		"template", "tempo", "grafana/tempo",
		"--version", "1.24.4",
		"--namespace", "observability-traces",
		"-f", "values/profiles/traces.yaml",
	}
	cmd := exec.Command(helm, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("tempo helm template failed: %w\n%s", err, stderr.String())
	}
	kubeconform, err := findKubeconform()
	if err != nil {
		if opts.StrictTools {
			return err
		}
		fmt.Println("kubeconform unavailable; skipped tempo rendered manifest validation")
	} else if err := runKubeconform(kubeconform, out); err != nil {
		return fmt.Errorf("tempo rendered manifest validation failed: %w", err)
	}

	if err := checkTracePlaceholders(); err != nil {
		return err
	}
	fmt.Println("traces validation ok")
	return nil
}

func checkTracePlaceholders() error {
	files := []string{
		"values/profiles/traces.yaml",
		"values/profiles/traces-prometheus.yaml",
		"values/overrides/single-cluster-traces.yaml",
		"examples/opentelemetry/traces-instrumentation.yaml",
		"examples/phase3-smoke/trace-generator.yaml",
	}
	forbidden := []string{
		"api.honeycomb.io",
		"ingest.",
		"arn:" + "aws",
		"tempo-chunks",
		"trace-tenant",
		"tenant_id",
	}
	var hits []string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		for _, token := range forbidden {
			if strings.Contains(string(data), token) {
				hits = append(hits, fmt.Sprintf("%s contains implementation-owned trace value %q", file, token))
			}
		}
	}
	if len(hits) > 0 {
		return fmt.Errorf("traces placeholder check failed:\n%s", strings.Join(hits, "\n"))
	}
	return nil
}
