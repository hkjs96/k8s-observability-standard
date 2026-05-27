package validate

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Logs(opts Options) error {
	helm, err := findHelm()
	if err != nil {
		if opts.StrictTools {
			return err
		}
		fmt.Println("helm unavailable; skipped logs helm template")
		return checkLogLabelCardinality()
	}

	renders := []struct {
		name    string
		chart   string
		version string
		values  string
	}{
		{name: "loki", chart: "grafana/loki", version: "7.0.0", values: "values/profiles/logs.yaml"},
		{name: "alloy", chart: "grafana/alloy", version: "1.8.2", values: "values/profiles/logs-alloy.yaml"},
	}

	for _, render := range renders {
		args := []string{
			"template", render.name, render.chart,
			"--version", render.version,
			"--namespace", "observability-logs",
			"-f", render.values,
		}
		out, err := exec.Command(helm, args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s helm template failed: %w\n%s", render.name, err, string(out))
		}
		kubeconform, err := findKubeconform()
		if err != nil {
			if opts.StrictTools {
				return err
			}
			fmt.Printf("kubeconform unavailable; skipped %s rendered manifest validation\n", render.name)
		} else if err := runKubeconform(kubeconform, out); err != nil {
			return fmt.Errorf("%s rendered manifest validation failed: %w", render.name, err)
		}
	}

	if err := checkLogLabelCardinality(); err != nil {
		return err
	}
	fmt.Println("logs validation ok")
	return nil
}

func checkLogLabelCardinality() error {
	files := []string{
		"values/profiles/logs.yaml",
		"values/profiles/logs-alloy.yaml",
	}
	forbidden := []string{
		"__meta_kubernetes_pod_uid",
		"__meta_kubernetes_pod_container_id",
		"__meta_kubernetes_pod_annotation_kubectl_kubernetes_io_restartedAt",
		"pod_uid",
		"container_id",
		"filename",
	}
	var hits []string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		for _, token := range forbidden {
			if strings.Contains(string(data), token) {
				hits = append(hits, fmt.Sprintf("%s contains high-cardinality log label source %q", file, token))
			}
		}
	}
	if len(hits) > 0 {
		return fmt.Errorf("logs label cardinality check failed:\n%s", strings.Join(hits, "\n"))
	}
	return nil
}
