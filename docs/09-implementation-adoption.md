# Implementation Adoption Guide

This guide describes how an implementation repository should consume this
standard without moving runtime-specific values into this repository.

## Repository Split

This standard repository owns:

- common Helm values
- profile values
- environment and sizing defaults
- Argo CD templates
- dashboard and rule baselines
- validation scripts and documentation

The implementation repository owns:

- final cluster and environment labels
- secret manifests or external secret references
- storage class, retention, and replica overrides
- alert receiver configuration
- ingress or private access configuration
- implementation-specific GitOps repository URL and branch
- runbook URLs and team ownership labels

## Adoption Flow

1. Copy or reference the Argo CD Application pattern from `argocd/`.
2. Keep the upstream values layering order:

```yaml
helm:
  valueFiles:
    - $values/values/common/kube-prometheus-stack.yaml
    - $values/values/profiles/basic.yaml
    - $values/values/env/prod.yaml
    - $values/values/sizing/medium.yaml
    - $values/values/overrides/implementation-specific.yaml
```

3. Create the `monitoring` namespace or adopt an existing platform namespace.
4. Create the `grafana-admin` Secret in the target namespace.
5. Add implementation-owned alert receiver configuration.
6. Render and validate before enabling automated Argo CD sync.

## Required Implementation Inputs

- Cluster identifier.
- Environment name.
- Storage class and retention policy.
- Grafana access path.
- Alert receiver and routing policy.
- Any ServiceMonitor or PodMonitor selectors that differ from the standard.
- Runbook URLs for alerts.

## Validation Before Sync

Run these checks from the implementation repository after wiring the standard:

```powershell
go run ./cmd/obsctl validate --strict-tools
helm template kube-prometheus-stack prometheus-community/kube-prometheus-stack `
  --version 85.0.2 `
  --namespace monitoring `
  -f values/common/kube-prometheus-stack.yaml `
  -f values/profiles/basic.yaml `
  -f values/env/prod.yaml `
  -f values/sizing/medium.yaml `
  -f values/overrides/implementation-specific.yaml
```

Use implementation-owned values for the final override file. Do not add those
runtime values to this repository.

## Handover Evidence

Before handover, capture:

- chart version
- values layer list
- validation command output
- Basic smoke evidence, when a disposable cluster test was run
- Grafana access method
- alert receiver owner
- storage and retention assumptions
- rollback method

Use `templates/basic-smoke-evidence.template.md` for smoke-test evidence and
`templates/operation-handover.template.md` as the handover skeleton.
