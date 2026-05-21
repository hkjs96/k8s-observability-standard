---
id: basic-validation
type: check
applies_to:
  - "values/**"
  - "argocd/**"
  - "rules/**"
  - "policies/**"
commands:
  - "go run ./cmd/obsctl validate --strict-tools basic"
---

# Basic Validation

Run:

```powershell
go run ./cmd/obsctl validate --strict-tools basic
```

Completion criteria:

- YAML parses successfully.
- Basic Helm values render with the pinned `kube-prometheus-stack` chart when Helm is available.
- Rendered manifests pass kubeconform when kubeconform is available.
- Helm lint passes when Helm is available.
- Missing optional external tools are reported clearly.
