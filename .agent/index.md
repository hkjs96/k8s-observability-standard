---
id: agent-index
type: index
required: true
applies_to:
  - "**/*"
---

# Agent Instruction Index

This directory is the primary execution context for AI agents. `docs/` is
human-facing reference material; use it for rationale and background after the
required `.agent/` files are read.

## Always Read

1. `.agent/context/product.md`
2. `.agent/context/repository-map.md`
3. `.agent/rules/repository-boundary.md`
4. `.agent/rules/sensitive-values.md`

## Path Routing

| Changed path | Required files |
| --- | --- |
| `values/**` | `.agent/rules/profiles.md`, `.agent/workflows/update-values.md`, `.agent/checks/basic-validation.md` |
| `argocd/**` | `.agent/rules/argocd.md`, `.agent/workflows/update-argocd-app.md`, `.agent/checks/argocd-validation.md` |
| `rules/**` | `.agent/workflows/add-prometheus-rule.md`, `.agent/checks/prometheus-rule-validation.md` |
| `dashboards/**` | `.agent/workflows/add-dashboard.md`, `.agent/checks/sensitive-values.md` |
| `policies/**` | `.agent/rules/repository-boundary.md`, `.agent/checks/basic-validation.md` |
| `docs/10-phase-3-plan.md` | `.agent/workflows/phase-3-profile-planning.md`, `.agent/checks/phase-3-profile-validation.md` |
| `docs/09-implementation-adoption.md`, `templates/implementation-adoption-checklist.template.md` | `.agent/workflows/implementation-adoption.md`, `.agent/checks/implementation-adoption-validation.md` |
| Phase 3 profile values, rules, dashboards, or policies | `.agent/workflows/phase-3-profile-planning.md`, `.agent/checks/phase-3-profile-validation.md` |
| `docs/**` | Use `.agent/rules/repository-boundary.md`; keep docs human-facing. |
| `templates/**` | `.agent/workflows/implementation-adoption.md`, `.agent/checks/implementation-adoption-validation.md` |
| `.agent/**` | Keep front matter accurate and update this index if routing changes. |
| `cmd/**`, `internal/**` | Run `go test ./...` and `go run ./cmd/obsctl validate`. |

## Completion Rule

Before reporting completion, run every check listed in the relevant rule or
workflow front matter unless the required tool is unavailable. If a tool is
unavailable, report exactly what was skipped and why.
