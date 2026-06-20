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
	registry := map[string]func() error{
		"yaml":       func() error { return YAML(opts) },
		"charts":     Charts,
		"basic":      func() error { return Basic(opts) },
		"logs":       func() error { return Logs(opts) },
		"traces":     func() error { return Traces(opts) },
		"slo":        func() error { return SLO(opts) },
		"argocd":     ArgoCD,
		"prometheus": func() error { return PrometheusRules(opts) },
		"sensitive":  SensitiveValues,
	}
	build := func(names ...string) []Check {
		checks := make([]Check, 0, len(names))
		for _, name := range names {
			checks = append(checks, Check{Name: name, Fn: registry[name]})
		}
		return checks
	}

	switch target {
	case "", "all":
		return build("yaml", "charts", "basic", "logs", "traces", "slo", "argocd", "prometheus", "sensitive"), nil
	case "profile/basic":
		return build("yaml", "basic", "prometheus", "sensitive"), nil
	case "profile/logs":
		return build("logs"), nil
	case "profile/traces":
		return build("traces"), nil
	case "profile/slo":
		return build("slo"), nil
	}
	if _, ok := registry[target]; ok {
		return build(target), nil
	}
	return nil, fmt.Errorf("unknown validation target %q", target)
}
