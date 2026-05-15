---
id: argocd-validation
type: check
applies_to:
  - "argocd/**"
commands:
  - "go run ./cmd/obsctl validate argocd"
---

# Argo CD Validation

Run:

```powershell
go run ./cmd/obsctl validate argocd
```

Completion criteria:

- Argo CD YAML parses.
- Applications do not use the `default` project.
- AppProjects do not use wildcard source repositories or wildcard destination namespaces.
- Helm Applications use `valueFiles`.
