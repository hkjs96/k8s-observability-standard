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

	want := []string{"yaml", "basic", "argocd", "prometheus", "sensitive"}
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

func TestSelectChecksRejectsUnknownTarget(t *testing.T) {
	_, err := selectChecks("unknown", Options{})
	if err == nil {
		t.Fatal("selectChecks() error = nil, want unknown target error")
	}
	if !strings.Contains(err.Error(), "unknown validation target") {
		t.Fatalf("selectChecks() error = %v, want unknown target detail", err)
	}
}
