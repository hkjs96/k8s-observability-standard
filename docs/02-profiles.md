# Profiles

Profiles are additive. Each profile must have a clear ownership boundary and
must not force the full observability stack on every implementation.

## Basic

Purpose: first monitoring baseline.

Includes:

- `kube-prometheus-stack`
- embedded Grafana
- default Kubernetes dashboards plus curated dashboard candidates
- PrometheusRule, ServiceMonitor, and PodMonitor conventions
- Alertmanager null route until the implementation repo supplies real routes

Excludes:

- Loki and log collectors
- Tempo and trace receivers
- SLO generators
- Mimir
- eBPF profiling

## Standard Logs

Purpose: add centralized log storage and query.

Includes later:

- Loki backend
- Alloy logs collector
- low-cardinality label policy
- Java multiline parser sample
- log push/query smoke tests

## Advanced Traces

Purpose: add distributed tracing and trace-to-logs workflows.

Includes later:

- Tempo backend
- OTLP receiver service
- OpenTelemetry Operator instrumentation samples
- Grafana trace-to-logs settings

## SLO

Purpose: define SLO-as-code workflows and error budget operations.

Includes later:

- Sloth validation and rule generation
- availability and latency SLO samples
- generated PrometheusRule review flow
- Pyrra decision criteria if a dedicated SLO UI is required

## Optional

Optional capabilities are evaluated only when implementation requirements justify
their cost and operational complexity:

- Mimir for central or long-term metrics
- Perses for dashboard-as-code PoC
- OpenObserve or Coroot as later UI/diagnostic candidates
- profiling or eBPF collectors after security review
