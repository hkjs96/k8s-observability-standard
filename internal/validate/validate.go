package validate

import (
	"fmt"
)

type Check struct {
	Name string
	Fn   func() error
}

func Run(target string) error {
	checks, err := selectChecks(target)
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

func selectChecks(target string) ([]Check, error) {
	all := []Check{
		{Name: "yaml", Fn: YAML},
		{Name: "basic", Fn: Basic},
		{Name: "argocd", Fn: ArgoCD},
		{Name: "prometheus", Fn: PrometheusRules},
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
