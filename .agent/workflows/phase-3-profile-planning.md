---
id: phase-3-profile-planning
type: workflow
required_for:
  - "docs/10-phase-3-plan.md"
  - "values/**"
  - "rules/**"
  - "dashboards/**"
  - "policies/**"
checks:
  - go run ./cmd/obsctl validate --strict-tools
  - go run ./cmd/obsctl validate profile basic --strict-tools
  - go run ./cmd/obsctl validate sensitive
---

# Phase 3 Profile Planning Workflow

1. Confirm the task explicitly expands beyond Phase 0-2 Basic.
2. Read `docs/10-phase-3-plan.md` and `.agent/rules/profiles.md`.
3. Keep each change scoped to one profile: logs, tracing, SLO, or a documented optional profile.
4. Do not add implementation-owned runtime values, secrets, endpoints, account IDs, bucket names, or identity bindings.
5. Add profile-scoped validation before publishing new profile resources.
6. Keep Basic profile validation green.
7. Run `go run ./cmd/obsctl validate --strict-tools`.
8. Run `go run ./cmd/obsctl validate profile basic --strict-tools`.
9. Run `go run ./cmd/obsctl validate sensitive`.
