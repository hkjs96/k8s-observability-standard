# Deployment Guide

This guide describes the Basic profile flow. The implementation repository owns
final environment values and secret creation.

## Prerequisites

- Kubernetes cluster with permission to install Prometheus Operator CRDs.
- Argo CD installed.
- Helm available for local rendering.
- `monitoring` namespace created by `argocd/applications/00-namespaces.yaml` or
  by an equivalent platform baseline.
- Grafana admin Secret created by the implementation repository:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: grafana-admin
  namespace: monitoring
type: Opaque
stringData:
  admin-user: REPLACE_IN_IMPLEMENTATION_REPO
  admin-password: REPLACE_IN_IMPLEMENTATION_REPO
```

Do not commit real secret values to this repository.

## Argo CD Order

Use sync waves:

| Wave | Resource |
| --- | --- |
| -20 | Namespace and Pod Security labels |
| 10 | `kube-prometheus-stack` |
| 50 | Dashboard and alert provisioning samples |
| 90 | Smoke tests |

CRDs and CRs must not be forced into the same lifecycle without an upgrade plan.

## Basic Smoke Tests

For a disposable Amazon Linux 2023 k3s runtime check, use
`docs/13-ec2-k3s-basic-smoke-runbook.md` and capture results with
`templates/basic-smoke-evidence.template.md`.

- Grafana datasource health succeeds.
- Prometheus targets are up.
- Alertmanager is reachable from Prometheus.
- Default and custom PrometheusRule resources are loaded.
- Dashboard ConfigMaps with `grafana_dashboard=1` are discovered by the sidecar.
