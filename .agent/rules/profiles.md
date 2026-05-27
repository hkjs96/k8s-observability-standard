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

- Phase 0-2 may modify the Basic profile only.
- Basic uses `kube-prometheus-stack` and embedded Grafana.
- Do not add Loki, Alloy, Tempo, Mimir, Sloth, Pyrra, or Perses unless the task explicitly expands scope.
- Keep upstream chart versions pinned in `charts-lock/chart-versions.yaml`.
- Keep values layered as common, profile, environment, sizing, then implementation override.
- Use `obsctl validate profile <name>` for profile-scoped validation when the target exists.
