package validate

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// sloEscalationContact matches committed escalation contacts (email addresses)
// that belong in an implementation repository, not in the standard SLO samples.
var sloEscalationContact = regexp.MustCompile(`[\w.+-]+@[\w-]+\.[\w.-]+`)

func SLO(opts Options) error {
	required := []string{
		"examples/slo/example-availability.slo.yaml",
		"examples/phase3-smoke/slo-metrics-generator.yaml",
		"dashboards/grafana/slo-overview.yaml",
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
		"examples/phase3-smoke/slo-metrics-generator.yaml",
		"rules/slo/example-availability.yaml",
		"templates/error-budget-review.template.md",
	}
	forbidden := []string{
		"pagerduty",
		"slack.com",
		"opsgenie",
		"arn:aws",
	}
	var hits []string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		text := string(data)
		lower := strings.ToLower(text)
		for _, token := range forbidden {
			if strings.Contains(lower, token) {
				hits = append(hits, fmt.Sprintf("%s contains implementation-owned SLO value %q", file, token))
			}
		}
		if contact := sloEscalationContact.FindString(text); contact != "" {
			hits = append(hits, fmt.Sprintf("%s contains implementation-owned escalation contact %q", file, contact))
		}
	}
	if len(hits) > 0 {
		return fmt.Errorf("slo placeholder check failed:\n%s", strings.Join(hits, "\n"))
	}
	return nil
}
