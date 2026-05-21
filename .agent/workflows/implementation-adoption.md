---
id: implementation-adoption
type: workflow
required_for:
  - "docs/09-implementation-adoption.md"
  - "templates/implementation-adoption-checklist.template.md"
  - "examples/**"
checks:
  - go run ./cmd/obsctl validate --strict-tools
  - go run ./cmd/obsctl validate sensitive
---

# Implementation Adoption Workflow

1. Confirm the change is guidance or a fictional example, not a runtime deployment.
2. Read `.agent/rules/repository-boundary.md` and `.agent/rules/sensitive-values.md`.
3. Keep real cluster, account, domain, bucket, webhook, and identity values out of this repository.
4. Keep final overrides and secrets owned by implementation repositories.
5. Update `templates/implementation-adoption-checklist.template.md` when adoption evidence changes.
6. Run `go run ./cmd/obsctl validate --strict-tools`.
7. Run `go run ./cmd/obsctl validate sensitive`.
