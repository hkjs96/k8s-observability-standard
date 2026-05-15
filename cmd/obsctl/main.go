package main

import (
	"fmt"
	"os"

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
		target := "all"
		if len(args) > 1 {
			target = args[1]
		}
		return validate.Run(target)
	case "help", "-h", "--help":
		return usage()
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage() error {
	fmt.Println(`obsctl validates the Kubernetes observability standard repository.

Usage:
  obsctl validate [all|basic|yaml|sensitive|argocd|prometheus]

Default:
  obsctl validate all`)
	return nil
}
