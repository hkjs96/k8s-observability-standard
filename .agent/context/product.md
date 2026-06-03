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

Implemented scope:

- standard repository scaffold
- human-facing standard documentation
- Basic profile based on `kube-prometheus-stack`
- Phase 3 profiles: Standard Logs (Loki + Alloy), Advanced Traces (Tempo), and
  SLO-as-code samples
- Argo CD templates using pinned chart versions and Git-hosted values
- Go validation CLI (`cmd/obsctl`) and templates

Still out of scope unless a task explicitly expands it:

- Mimir for central or long-term metrics
- Sloth or Pyrra as the standardized SLO generator or UI
- Perses dashboard-as-code
- eBPF or profiling collectors and other optional UIs

When extending a profile, follow `.agent/workflows/phase-3-profile-planning.md`
and keep `charts-lock/chart-versions.yaml` as the single source of truth for
pinned chart versions.
