---
id: update-values
type: workflow
required_for:
  - "values/**"
checks:
  - go run ./cmd/obsctl validate basic
  - go run ./cmd/obsctl validate sensitive
---

# Update Values Workflow

1. Read `charts-lock/chart-versions.yaml` before changing chart-specific values.
2. Keep changes in the smallest relevant layer: common, profile, environment, or sizing.
3. Do not put implementation-owned values in this repository.
4. Keep Grafana credentials as existing Secret references.
5. Render the Basic profile with `go run ./cmd/obsctl validate basic`.
