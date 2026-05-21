package main

import "testing"

func TestParseValidateArgsDefaults(t *testing.T) {
	target, opts, err := parseValidateArgs(nil)
	if err != nil {
		t.Fatal(err)
	}
	if target != "all" {
		t.Fatalf("target = %q, want all", target)
	}
	if opts.StrictTools {
		t.Fatal("StrictTools = true, want false")
	}
}

func TestParseValidateArgsStrictToolsAndTarget(t *testing.T) {
	target, opts, err := parseValidateArgs([]string{"--strict-tools", "prometheus"})
	if err != nil {
		t.Fatal(err)
	}
	if target != "prometheus" {
		t.Fatalf("target = %q, want prometheus", target)
	}
	if !opts.StrictTools {
		t.Fatal("StrictTools = false, want true")
	}
}

func TestParseValidateArgsRejectsUnknown(t *testing.T) {
	if _, _, err := parseValidateArgs([]string{"--bad"}); err == nil {
		t.Fatal("parseValidateArgs() error = nil, want error")
	}
}
