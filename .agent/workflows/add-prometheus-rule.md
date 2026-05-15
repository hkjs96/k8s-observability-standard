---
id: add-prometheus-rule
type: workflow
required_for:
  - "rules/prometheus/**"
  - "rules/alerting/**"
checks:
  - go run ./cmd/obsctl validate prometheus
  - go run ./cmd/obsctl validate sensitive
---

# Add Prometheus Rule Workflow

1. Add or update the Prometheus Operator `PrometheusRule` resource.
2. Add or update a matching `.promtool.yaml` file containing only rule groups.
3. Keep runbook URLs fictional unless they belong to implementation repos.
4. Run `go run ./cmd/obsctl validate prometheus`.
5. Run `go run ./cmd/obsctl validate sensitive`.
