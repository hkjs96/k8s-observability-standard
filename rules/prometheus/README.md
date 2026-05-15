# Prometheus Rules

Deployable samples use the Prometheus Operator `PrometheusRule` resource.
Files with `.promtool.yaml` mirror only the `groups` section so `promtool check
rules` can validate the expressions without Kubernetes CRD wrapping.

Keep generated or implementation-specific rule files in the implementation
repository. This directory is for reusable baseline examples.
