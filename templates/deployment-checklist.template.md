# Deployment Checklist

## Pre-Deployment

- Chart version matches `charts-lock/chart-versions.yaml`.
- Values are layered through `valueFiles`.
- Grafana admin Secret exists in the target namespace.
- No real secrets or environment-specific values are committed to this standard repository.
- Argo CD AppProject does not use wildcard source repositories or destination namespaces.

## Static Validation

- YAML parse completed.
- `helm lint` completed.
- `helm template` completed.
- `kubeconform` or equivalent manifest validation completed.
- `promtool check rules` completed for `.promtool.yaml` files.

## Smoke Tests

- Smoke evidence template copied to the implementation repository.
- Prometheus pods ready.
- Alertmanager pods ready.
- Grafana pods ready.
- Prometheus targets are up.
- Grafana datasource health succeeds.
- Dashboard sidecar discovers `grafana_dashboard=1` ConfigMaps.
- Test alert route is verified in the implementation repository.

## Rollback

- Previous chart version:
- Previous values commit:
- Rollback command:
- Known data retention impact:
