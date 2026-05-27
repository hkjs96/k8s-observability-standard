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

func TestParseValidateArgsProfileBasic(t *testing.T) {
	target, opts, err := parseValidateArgs([]string{"profile", "basic", "--strict-tools"})
	if err != nil {
		t.Fatal(err)
	}
	if target != "profile/basic" {
		t.Fatalf("target = %q, want profile/basic", target)
	}
	if !opts.StrictTools {
		t.Fatal("StrictTools = false, want true")
	}
}

func TestParseValidateArgsRejectsUnknownProfile(t *testing.T) {
	if _, _, err := parseValidateArgs([]string{"profile", "mimir"}); err == nil {
		t.Fatal("parseValidateArgs() error = nil, want error")
	}
}

func TestParseValidateArgsRejectsUnknown(t *testing.T) {
	if _, _, err := parseValidateArgs([]string{"--bad"}); err == nil {
		t.Fatal("parseValidateArgs() error = nil, want error")
	}
}
