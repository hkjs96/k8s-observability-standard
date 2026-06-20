# Grafana Dashboard Candidates

The Basic profile starts with the upstream dashboards shipped by
`kube-prometheus-stack`. Additional dashboards from prior PoC work must be
reviewed and sanitized before being copied here.

## Added

- `cluster-overview.yaml` - CPU/memory by namespace, node CPU utilisation
- `namespace-overview.yaml` - per-namespace CPU/memory/pod status with namespace variable
- `observability-datasources.yaml` - Loki and Tempo datasource provisioning for profile smoke tests
- `logs-overview.yaml` - Loki log volume, error-like logs, and recent logs
- `traces-overview.yaml` - Tempo target health and trace ingestion metrics
- `slo-overview.yaml` - example availability SLO ratio, request rate, and burn-rate panels
- `profile-alerting.yaml` - Grafana-managed alert rules for logs/traces profile checks
- `profile-notifiers.yaml` - example Grafana contact point and notification policy

## Remaining Candidates For Phase 2 Review

- Node overview
- Workload resources
- Pod status timeline
- HPA overview
- Service overview
- Prometheus self-monitoring

## Adoption Rules

- Dashboard ConfigMaps must use `grafana_dashboard=1`.
- Datasource ConfigMaps must use `grafana_datasource=1`.
- Grafana-managed alert ConfigMaps must use `grafana_alert=1`.
- Grafana contact point and notification policy ConfigMaps must use `grafana_alert=1`.
- Folder placement should use the `grafana_folder` annotation.
- Dashboards must not contain real cluster names, datasource UIDs from a live
  environment, hostnames, account IDs, or customer labels.
- Keep Grafana as the Basic operating UI. Perses is evaluated in a later PoC and
  is not part of Phase 0-2.

## ConfigMap Shape

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboard-cluster-overview
  namespace: monitoring
  labels:
    grafana_dashboard: "1"
  annotations:
    grafana_folder: Kubernetes
data:
  cluster-overview.json: |
    {}
```
