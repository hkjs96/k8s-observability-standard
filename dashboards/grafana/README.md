# Grafana Dashboard Candidates

The Basic profile starts with the upstream dashboards shipped by
`kube-prometheus-stack`. Additional dashboards from prior PoC work must be
reviewed and sanitized before being copied here.

## Candidate Set For Phase 2 Review

- Cluster overview
- Node overview
- Namespace overview
- Workload resources
- Pod status timeline
- HPA overview
- Service overview
- Prometheus self-monitoring

## Adoption Rules

- Dashboard ConfigMaps must use `grafana_dashboard=1`.
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
