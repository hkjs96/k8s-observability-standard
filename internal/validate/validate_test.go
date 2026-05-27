package validate

import (
	"reflect"
	"strings"
	"testing"
)

func TestSelectChecksAllTargets(t *testing.T) {
	checks, err := selectChecks("", Options{})
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	for _, check := range checks {
		names = append(names, check.Name)
	}

	want := []string{"yaml", "basic", "logs", "traces", "slo", "argocd", "prometheus", "sensitive"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("selectChecks(\"\") = %v, want %v", names, want)
	}
}

func TestSelectChecksSingleTarget(t *testing.T) {
	checks, err := selectChecks("prometheus", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 1 || checks[0].Name != "prometheus" {
		t.Fatalf("selectChecks(\"prometheus\") = %v, want one prometheus check", checks)
	}
}

func TestSelectChecksProfileBasic(t *testing.T) {
	checks, err := selectChecks("profile/basic", Options{})
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	for _, check := range checks {
		names = append(names, check.Name)
	}

	want := []string{"yaml", "basic", "prometheus", "sensitive"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("selectChecks(\"profile/basic\") = %v, want %v", names, want)
	}
}

func TestSelectChecksProfileLogs(t *testing.T) {
	checks, err := selectChecks("profile/logs", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 1 || checks[0].Name != "logs" {
		t.Fatalf("selectChecks(\"profile/logs\") = %v, want one logs check", checks)
	}
}

func TestSelectChecksProfileSLO(t *testing.T) {
	checks, err := selectChecks("profile/slo", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 1 || checks[0].Name != "slo" {
		t.Fatalf("selectChecks(\"profile/slo\") = %v, want one slo check", checks)
	}
}

func TestSelectChecksProfileTraces(t *testing.T) {
	checks, err := selectChecks("profile/traces", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 1 || checks[0].Name != "traces" {
		t.Fatalf("selectChecks(\"profile/traces\") = %v, want one traces check", checks)
	}
}

func TestSelectChecksRejectsUnknownTarget(t *testing.T) {
	_, err := selectChecks("unknown", Options{})
	if err == nil {
		t.Fatal("selectChecks() error = nil, want unknown target error")
	}
	if !strings.Contains(err.Error(), "unknown validation target") {
		t.Fatalf("selectChecks() error = %v, want unknown target detail", err)
	}
}
