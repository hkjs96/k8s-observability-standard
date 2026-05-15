---
id: argocd
type: rule
required: false
applies_to:
  - "argocd/**"
checks:
  - go run ./cmd/obsctl validate argocd
  - go run ./cmd/obsctl validate sensitive
---

# Argo CD Rules

- Use explicit AppProject names.
- Do not use the `default` project for production templates.
- Do not use wildcard `sourceRepos`.
- Do not use wildcard destination namespaces.
- Prefer Helm `valueFiles` over inline `values` or `valuesObject`.
- Use Argo CD multiple sources when combining an upstream chart with Git-hosted values.
- Keep sync waves explicit for namespace and core platform resources.
