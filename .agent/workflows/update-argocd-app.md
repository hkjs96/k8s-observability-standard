---
id: update-argocd-app
type: workflow
required_for:
  - "argocd/**"
checks:
  - go run ./cmd/obsctl validate argocd
  - go run ./cmd/obsctl validate charts
  - go run ./cmd/obsctl validate sensitive
---

# Update Argo CD Application Workflow

1. Read `.agent/rules/argocd.md`.
2. Keep AppProject boundaries explicit.
3. Use Helm `valueFiles` and pinned chart versions that match `charts-lock/chart-versions.yaml`.
4. Do not use wildcard source repositories or wildcard destination namespaces.
5. Run `go run ./cmd/obsctl validate argocd`.
6. Run `go run ./cmd/obsctl validate charts` to confirm each `targetRevision` matches the chart lock.
