package validate

import (
	"fmt"
)

type Check struct {
	Name string
	Fn   func() error
}

type Options struct {
	StrictTools bool
}

func Run(target string, opts Options) error {
	checks, err := selectChecks(target, opts)
	if err != nil {
		return err
	}

	for _, check := range checks {
		fmt.Printf("==> %s\n", check.Name)
		if err := check.Fn(); err != nil {
			return err
		}
	}

	fmt.Println("validation ok")
	return nil
}

func selectChecks(target string, opts Options) ([]Check, error) {
	all := []Check{
		{Name: "yaml", Fn: func() error { return YAML(opts) }},
		{Name: "basic", Fn: func() error { return Basic(opts) }},
		{Name: "argocd", Fn: ArgoCD},
		{Name: "prometheus", Fn: func() error { return PrometheusRules(opts) }},
		{Name: "sensitive", Fn: SensitiveValues},
	}

	if target == "" || target == "all" {
		return all, nil
	}
	for _, check := range all {
		if check.Name == target {
			return []Check{check}, nil
		}
	}
	return nil, fmt.Errorf("unknown validation target %q", target)
}
