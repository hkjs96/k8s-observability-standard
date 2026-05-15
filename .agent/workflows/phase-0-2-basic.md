---
id: phase-0-2-basic
type: workflow
required_for:
  - "values/**"
  - "argocd/**"
  - "rules/**"
  - "policies/**"
checks:
  - go run ./cmd/obsctl validate basic
  - go run ./cmd/obsctl validate sensitive
---

# Phase 0-2 Basic Workflow

1. Confirm the task is still within Phase 0-2.
2. Read `.agent/rules/profiles.md`.
3. Update only Basic-related values, Argo CD templates, rules, policies, docs, or templates.
4. Do not introduce logs, traces, SLO, Mimir, Perses, or customer-specific runtime values.
5. Run `go run ./cmd/obsctl validate basic`.
6. Run `go run ./cmd/obsctl validate sensitive`.
