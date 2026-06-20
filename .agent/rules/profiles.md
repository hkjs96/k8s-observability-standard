---
id: profiles
type: rule
required: false
applies_to:
  - "values/**"
  - "charts-lock/**"
checks:
  - go run ./cmd/obsctl validate profile basic
  - go run ./cmd/obsctl validate sensitive
---

# Profile Rules

- Implemented profiles: Basic (`kube-prometheus-stack` + embedded Grafana),
  Standard Logs (Loki + Alloy), Advanced Traces (Tempo), and SLO-as-code samples.
- Basic uses `kube-prometheus-stack` and embedded Grafana.
- Do not add Mimir, Sloth, Pyrra, or Perses unless the task explicitly expands scope.
- Keep upstream chart versions pinned in `charts-lock/chart-versions.yaml`; it is
  the single source of truth that `obsctl validate charts` enforces.
- Keep values layered as common, profile, environment, sizing, then implementation override.
- Use `obsctl validate profile <name>` for profile-scoped validation when the target exists.
