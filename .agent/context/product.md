---
id: product-context
type: context
required: true
applies_to:
  - "**/*"
---

# Product Context

This repository is a reusable Kubernetes Observability v2 standard repository.
It is not a live operations repository and not a central multi-customer repo.

Phase 0-2 scope:

- standard repository scaffold
- human-facing standard documentation
- Basic profile based on `kube-prometheus-stack`
- Argo CD templates using pinned chart versions and Git-hosted values
- validation scripts and templates

Later profiles such as Loki, Alloy, Tempo, Mimir, Sloth, Pyrra, and Perses must
not be implemented unless the task explicitly expands scope.
