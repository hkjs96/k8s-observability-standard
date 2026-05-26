# Implementation Adoption Checklist

Use this checklist in an implementation repository before enabling automated
sync. Keep real runtime values outside this standard repository.

## Repository Wiring

- Standard repository source selected:
- Standard repository revision selected:
- Implementation repository override path selected:
- Argo CD Application references the expected value files:
- Argo CD AppProject avoids wildcard source repositories:
- Argo CD destination namespace is explicit:

## Required Runtime Inputs

- `monitoring` namespace exists or is created by platform baseline:
- `grafana-admin` Secret exists in the target namespace:
- Cluster label selected:
- Environment label selected:
- Storage class selected:
- Prometheus retention selected:
- Alert receiver owner selected:
- Grafana access method selected:

## Static Validation

- `go run ./cmd/obsctl validate --strict-tools` completed:
- YAML parse completed:
- Helm template completed:
- Helm lint completed:
- kubeconform completed:
- promtool completed:
- Sensitive value scan completed:

## Deployment Review

- Chart version matches `charts-lock/chart-versions.yaml`:
- Values layer order reviewed:
- Implementation override contains no plaintext credentials:
- Alert routes reviewed:
- Dashboard folder and labels reviewed:
- Rollback commit identified:

## Smoke Test

- Smoke evidence captured with `templates/basic-smoke-evidence.template.md`:
- Prometheus target health checked:
- kube-state-metrics data visible:
- node-exporter data visible:
- Grafana datasource healthy:
- Cluster overview dashboard loads:
- Namespace overview dashboard loads:
- Test alert route verified, if enabled:

## Handover

- Operation owner:
- Alert owner:
- Grafana access owner:
- Retention assumption:
- Storage assumption:
- Rollback method:
- Known limits:
