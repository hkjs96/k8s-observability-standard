---
id: repository-map
type: context
required: true
applies_to:
  - "**/*"
---

# Repository Map

- `.agent/`: mandatory agent routing, rules, workflows, and checks
- `AGENTS.md`: Codex and generic agent entrypoint
- `CLAUDE.md`: Claude Code entrypoint
- `docs/`: human-facing standard documents and rationale
- `charts-lock/`: pinned upstream chart versions
- `values/`: reusable Helm values layers
- `argocd/`: AppProject and Application templates
- `dashboards/`: dashboard candidates and sanitized dashboard assets
- `rules/`: PrometheusRule resources and promtool mirror files
- `policies/`: security and namespace policy samples
- `templates/`: implementation assessment, security review, deployment, and handover templates
- `examples/`: fictional implementation examples only
- `cmd/obsctl`: Go validation and smoke-test CLI
