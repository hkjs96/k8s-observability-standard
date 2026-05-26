package main

import (
	"fmt"
	"os"

	"github.com/example/k8s-observability/internal/smoke"
	"github.com/example/k8s-observability/internal/validate"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}

	switch args[0] {
	case "validate":
		target, opts, err := parseValidateArgs(args[1:])
		if err != nil {
			return err
		}
		return validate.Run(target, opts)
	case "smoke":
		return smoke.Run(args[1:])
	case "help", "-h", "--help":
		return usage()
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func parseValidateArgs(args []string) (string, validate.Options, error) {
	target := "all"
	var opts validate.Options
	for _, arg := range args {
		switch arg {
		case "--strict-tools":
			opts.StrictTools = true
		case "all", "basic", "yaml", "sensitive", "argocd", "prometheus":
			target = arg
		default:
			return "", opts, fmt.Errorf("unknown validate argument %q", arg)
		}
	}
	return target, opts, nil
}

func usage() error {
	fmt.Println(`obsctl validates the Kubernetes observability standard repository.

Usage:
  obsctl validate [--strict-tools] [all|basic|yaml|sensitive|argocd|prometheus]
  obsctl smoke ec2-k3s [launch|fetch-kubeconfig|terminate] [options]
  obsctl smoke k3s-basic install [options]

Default:
  obsctl validate all

Options:
  --strict-tools  fail when optional validation tools are unavailable

Run "obsctl smoke help" for disposable smoke helper usage.`)
	return nil
}
