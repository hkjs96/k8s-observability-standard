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
2. Add or update the `.promtool.yaml` mirror so its rule groups match the
   `PrometheusRule` `spec.groups` exactly; the mirror drops only the Kubernetes
   wrapper (apiVersion, kind, metadata).
3. Keep runbook URLs fictional unless they belong to implementation repos.
4. Run `go run ./cmd/obsctl validate prometheus`.
5. Run `go run ./cmd/obsctl validate sensitive`.
