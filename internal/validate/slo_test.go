package validate

import "testing"

func TestCheckSLOPlaceholdersAcceptsExampleArtifacts(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, root, "examples/slo/example-availability.slo.yaml", "service: example-service\n")
	mustWrite(t, root, "examples/phase3-smoke/slo-metrics-generator.yaml", "service: example-service\n")
	mustWrite(t, root, "rules/slo/example-availability.yaml", "runbook_url: https://runbooks.example.invalid/slo\n")
	mustWrite(t, root, "templates/error-budget-review.template.md", "Owner: REPLACE_IN_IMPLEMENTATION_REPO\n")
	withWorkingDir(t, root, func() {
		if err := checkSLOPlaceholders(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCheckSLOPlaceholdersRejectsEscalationValues(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, root, "examples/slo/example-availability.slo.yaml", "service: example-service\n")
	mustWrite(t, root, "examples/phase3-smoke/slo-metrics-generator.yaml", "service: example-service\n")
	mustWrite(t, root, "rules/slo/example-availability.yaml", "receiver: pagerduty\n")
	mustWrite(t, root, "templates/error-budget-review.template.md", "Owner: REPLACE_IN_IMPLEMENTATION_REPO\n")
	withWorkingDir(t, root, func() {
		if err := checkSLOPlaceholders(); err == nil {
			t.Fatal("checkSLOPlaceholders() error = nil, want implementation value error")
		}
	})
}
