# Phase 0-2 Completion Report

This report summarizes the current Basic profile baseline. It is human-facing
status documentation; mandatory agent rules remain under `.agent/`.

## Status

Phase 0-2 is complete enough to use as a reusable Basic observability standard
baseline. The repository now includes values, Argo CD templates, dashboards,
Prometheus rules, policy samples, implementation templates, examples, and local
validation tooling.

## Completed Outputs

- Basic `kube-prometheus-stack` profile with embedded Grafana.
- Common, profile, environment, sizing, and example override values.
- Pinned chart metadata in `charts-lock/chart-versions.yaml`.
- Argo CD AppProject and Application templates.
- Grafana dashboard ConfigMaps for cluster and namespace overview.
- Prometheus alerting and recording rules with matching `.promtool.yaml` mirrors.
- Security guardrail samples for network policy, namespace policy, and read-only RBAC.
- Validation CLI under `cmd/obsctl` and PowerShell fallback scripts.
- CI templates for GitHub Actions, GitLab CI, and Jenkins.

## Validation Evidence

Last local validation performed:

```powershell
go test ./...
powershell -ExecutionPolicy Bypass -File scripts\validate-all.ps1
```

Observed result:

- Go tests passed for validation and file walking packages.
- YAML parse check passed.
- Helm template and lint passed for chart version `85.0.2`.
- Rendered manifest validation passed with kubeconform.
- Argo CD template checks passed.
- `promtool check rules` passed for 2 mirror files with 5 total rules.
- Sensitive value scan passed.

## Deployment Readiness

The baseline is ready for implementation repository adoption when these runtime
inputs are supplied outside this repository:

- Grafana admin Secret named `grafana-admin` in the `monitoring` namespace.
- Final environment labels such as cluster and environment.
- Storage, retention, and sizing overrides for the target cluster.
- Alert receiver configuration and notification routes.
- Repository URL and branch values for the implementation GitOps source.

## Known Limits

- Phase 0-2 covers metrics and embedded Grafana only.
- Logs, tracing, SLO generation, long-term metrics, and alternate dashboard UI
  are documented as later profile work, not implemented here.
- CI templates are portable examples. Each implementation may need provider
  specific cache, network, and tool installation adjustments.
- `kubeconform` uses `-ignore-missing-schemas` because some CRDs do not have
  schemas available in the default schema catalog.

## Recommended Next Work

1. Validate the first implementation repository adoption against this standard.
2. Turn the Phase 3 plan into scoped logs, tracing, and SLO implementation tasks.
