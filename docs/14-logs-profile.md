# Logs Profile

The logs profile adds Loki for log storage and Grafana Alloy as the Kubernetes
log collector. Runtime storage, tenant, retention, and routing decisions remain
implementation-owned.

## Values

- `values/profiles/logs.yaml`: reusable Loki baseline.
- `values/profiles/logs-alloy.yaml`: reusable Alloy log collector baseline.
- `values/overrides/single-cluster-logs.yaml`: fictional single-cluster override.

The baseline renders Loki in `SingleBinary` mode with filesystem storage so the
standard can be validated locally without committing real object storage
details. Production implementations should decide whether to keep single-cluster
storage or use implementation-managed object storage.

## Sizing Guidance

- Small disposable checks: start with 50 GiB persistence and 7 days retention.
- Medium clusters: size retention from observed daily ingest and keep at least
  30 percent free headroom.
- Production clusters: use implementation-owned object storage, retention, and
  compaction settings.

Do not commit real bucket names, tenant IDs, cloud account IDs, or customer
endpoints in this repository.

## Label Policy

Keep log labels low-cardinality. The reusable Alloy config keeps:

- `namespace`
- `app`
- `container`

Do not promote pod UID, container ID, filename, request ID, user ID, or trace ID
to Loki labels. Keep those values in the log body.

## Validation

```powershell
go run ./cmd/obsctl validate profile logs --strict-tools
go run ./cmd/obsctl validate --strict-tools
```

The logs validator renders Loki and Alloy, runs kubeconform on rendered
manifests, and rejects high-cardinality log label sources.

## Grafana UI

Logs profile UI is managed through Grafana provisioning:

- `dashboards/grafana/observability-datasources.yaml`
- `dashboards/grafana/logs-overview.yaml`
- `dashboards/grafana/profile-alerting.yaml`

Use Grafana-managed alerts for dashboard-driven log conditions such as
error-like log lines. Keep platform-critical scrape and component availability
alerts in PrometheusRule resources.

## Smoke Queries

Run from an implementation repository after deployment:

```logql
{namespace="monitoring"}
{namespace="monitoring", app="kube-prometheus-stack"}
count_over_time({namespace="monitoring"}[5m])
```

Evidence to collect:

- Loki and Alloy pods are ready.
- Recent logs exist for the target namespace.
- Query latency is acceptable for the selected retention and storage mode.
- No high-cardinality labels appear in the Loki label browser.

## Rollback

Rollback is implementation-owned. At minimum:

- Disable the Alloy collector first to stop new ingestion.
- Preserve or snapshot Loki storage if logs are required for incident review.
- Remove the Loki release after evidence retention requirements are met.
