# Traces Profile

The traces profile adds Tempo for trace storage and OTLP ingestion. Real tenant,
storage, sampling, and exporter choices remain implementation-owned.

## Values

- `values/profiles/traces.yaml`: reusable Tempo baseline.
- `values/overrides/single-cluster-traces.yaml`: fictional single-cluster override.
- `examples/opentelemetry/traces-instrumentation.yaml`: OpenTelemetry Operator sample.

The baseline renders Tempo single binary mode with local trace storage so this
repository can validate manifests without committing object storage details.
The upstream chart is deprecated, so production adoption should revisit the
chart choice before rollout.

## OTLP Receiver Guidance

The reusable profile exposes:

- OTLP gRPC on `4317`
- OTLP HTTP on `4318`

Implementation repositories should route application exporters to the in-cluster
Tempo service or to an implementation-managed collector. Do not commit external
vendor ingestion URLs or tenant identifiers in this repository.

## Trace-To-Logs Guidance

Grafana trace-to-logs mappings should use stable labels shared with the logs
profile:

- `namespace`
- `app`
- `container`

Do not require pod UID, container ID, request ID, or user ID as labels. Keep
those values in span attributes or log bodies.

## Validation

```powershell
go run ./cmd/obsctl validate profile traces --strict-tools
go run ./cmd/obsctl validate --strict-tools
```

The traces validator renders Tempo, runs kubeconform on rendered manifests, and
rejects implementation-owned trace endpoint and tenant patterns.

## Grafana UI

Traces profile UI is managed through Grafana provisioning:

- `dashboards/grafana/observability-datasources.yaml`
- `dashboards/grafana/traces-overview.yaml`
- `dashboards/grafana/profile-alerting.yaml`

Use Grafana-managed alerts for trace and dashboard-driven checks. Keep
platform-critical scrape and component availability alerts in PrometheusRule
resources.

## Smoke Check

Run from an implementation repository after deployment:

```bash
otel-cli span \
  --service example-smoke \
  --name traces-smoke \
  --endpoint http://tempo.observability-traces.svc.cluster.local:4317
```

Evidence to collect:

- Tempo pod is ready.
- OTLP gRPC and HTTP ports are reachable from an application namespace.
- A smoke trace can be found in Grafana.
- Trace-to-logs links use low-cardinality labels.

## Rollback

Rollback is implementation-owned. At minimum:

- Stop application trace exporters or route them away from Tempo.
- Preserve Tempo storage if traces are required for incident review.
- Remove the Tempo release after retention requirements are met.
