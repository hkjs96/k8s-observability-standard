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
	basicProfile := []Check{
		{Name: "yaml", Fn: func() error { return YAML(opts) }},
		{Name: "basic", Fn: func() error { return Basic(opts) }},
		{Name: "prometheus", Fn: func() error { return PrometheusRules(opts) }},
		{Name: "sensitive", Fn: SensitiveValues},
	}
	all := []Check{
		{Name: "yaml", Fn: func() error { return YAML(opts) }},
		{Name: "basic", Fn: func() error { return Basic(opts) }},
		{Name: "logs", Fn: func() error { return Logs(opts) }},
		{Name: "argocd", Fn: ArgoCD},
		{Name: "prometheus", Fn: func() error { return PrometheusRules(opts) }},
		{Name: "sensitive", Fn: SensitiveValues},
	}

	if target == "" || target == "all" {
		return all, nil
	}
	switch target {
	case "profile/basic":
		return basicProfile, nil
	case "profile/logs":
		return []Check{{Name: "logs", Fn: func() error { return Logs(opts) }}}, nil
	case "profile/traces", "profile/slo":
		return nil, fmt.Errorf("validation target %q is not implemented because the profile is not present in this repository", target)
	}
	for _, check := range all {
		if check.Name == target {
			return []Check{check}, nil
		}
	}
	return nil, fmt.Errorf("unknown validation target %q", target)
}
