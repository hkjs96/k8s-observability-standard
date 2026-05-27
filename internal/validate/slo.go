package validate

import (
	"fmt"
	"os"
	"strings"
)

func SLO(opts Options) error {
	required := []string{
		"examples/slo/example-availability.slo.yaml",
		"rules/slo/example-availability.yaml",
		"rules/slo/example-availability.promtool.yaml",
		"templates/error-budget-review.template.md",
	}
	for _, file := range required {
		if _, err := os.Stat(file); err != nil {
			return fmt.Errorf("missing SLO artifact: %s", file)
		}
	}
	if err := PrometheusRules(opts); err != nil {
		return err
	}
	if err := checkSLOPlaceholders(); err != nil {
		return err
	}
	fmt.Println("slo validation ok")
	return nil
}

func checkSLOPlaceholders() error {
	files := []string{
		"examples/slo/example-availability.slo.yaml",
		"rules/slo/example-availability.yaml",
		"templates/error-budget-review.template.md",
	}
	forbidden := []string{
		"pagerduty",
		"slack.com",
		"@",
		"arn:" + "aws",
	}
	var hits []string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		for _, token := range forbidden {
			if strings.Contains(strings.ToLower(string(data)), token) {
				hits = append(hits, fmt.Sprintf("%s contains implementation-owned SLO value %q", file, token))
			}
		}
	}
	if len(hits) > 0 {
		return fmt.Errorf("slo placeholder check failed:\n%s", strings.Join(hits, "\n"))
	}
	return nil
}
