# Phase 3 Plan

Phase 3 expands beyond the Basic metrics baseline. The work should remain split
by profile so logs, tracing, and SLO features can be adopted independently.

> **Status:** The Standard Logs, Advanced Traces, and SLO profiles described
> below are now implemented and validation-ready — see `docs/14-logs-profile.md`,
> `docs/15-traces-profile.md`, and `docs/16-slo-profile.md`. The previously open
> decisions are now resolved per profile below. Work that extends beyond these
> three profiles (central metrics, the Sloth generation pipeline, a dedicated SLO
> UI, dashboard-as-code, and profiling) moves to `docs/18-phase-4-plan.md`. The
> shipped SLO sample is still a hand-written generated-style PrometheusRule (page
> and warning burn-rate alerts), not Sloth output, until that pipeline lands.

## Scope Principles

- Keep Basic metrics stable while adding new profiles.
- Add one backend or workflow at a time.
- Keep runtime values in implementation repositories.
- Add validation before publishing each new profile.
- Prefer upstream chart defaults unless the standard needs an explicit guardrail.

## Standard Logs Profile

Goal: add centralized log storage and query while keeping label cardinality
controlled.

Planned outputs:

- Loki values profile.
- Alloy log collector profile.
- Low-cardinality label policy.
- Java multiline parser sample.
- Log query smoke test.
- Local log generator smoke workload.
- Storage and retention guidance.

Validation additions:

- Helm render and kubeconform for log components.
- Static check for forbidden high-cardinality labels.
- Example LogQL smoke queries.

Resolved decisions:

- The standard ships a Loki `SingleBinary` filesystem baseline
  (`values/profiles/logs.yaml`). Central or multi-tenant backends are
  implementation-owned; long-term central metrics are tracked separately in
  `docs/18-phase-4-plan.md`.
- Object storage is optional for small or disposable tiers (filesystem) and
  expected at the production tier. Bucket names and credentials stay
  implementation-owned.
- Retention follows the sizing tiers in `docs/14-logs-profile.md`: 7 days for
  small checks, ingest-based with 30 percent headroom for medium clusters, and
  implementation-owned settings for production.

## Advanced Traces Profile

Goal: add distributed tracing and trace-to-logs workflows.

Planned outputs:

- Tempo values profile.
- OTLP receiver Service and routing guidance.
- OpenTelemetry Operator instrumentation samples.
- Grafana trace-to-logs settings.
- Trace smoke test.
- Local telemetrygen trace smoke workload.

Validation additions:

- Helm render and kubeconform for trace components.
- Static checks for OTLP endpoint placeholders.
- Example trace ingestion smoke command.

Resolved decisions:

- Small deployments use the Tempo single-binary local-storage baseline
  (`values/profiles/traces.yaml`); production object storage is
  implementation-owned. The upstream chart deprecation is flagged in
  `docs/15-traces-profile.md` and revisited in `docs/18-phase-4-plan.md`.
- The OpenTelemetry Operator is optional. The standard ships an instrumentation
  sample (`examples/opentelemetry/traces-instrumentation.yaml`) but does not
  require the operator.
- Default sampling guidance is parent-based head sampling: full sampling for
  dev and smoke, reduced by volume in production. The standard does not commit a
  fixed percentage, and tail sampling stays implementation-owned.

## SLO Profile

Goal: define SLO-as-code workflows and error budget operations.

Planned outputs:

- Sloth rule generation workflow.
- Availability and latency SLO examples.
- Alert routing guidance for burn-rate alerts.
- Error budget review template.
- Local SLO metrics generator and Grafana SLO dashboard.

Validation additions:

- SLO spec linting.
- Generated PrometheusRule mirror validation.
- promtool check for generated burn-rate rules.

Resolved decisions:

- The standard standardizes on Sloth for SLO-as-code generation, consistent
  with `docs/01-standard-decisions.md`. Pyrra stays an optional SLO UI evaluated
  in `docs/18-phase-4-plan.md`. Until the Sloth pipeline lands, the shipped
  sample is a hand-written generated-style PrometheusRule.
- SLO selection requires stable low-cardinality service labels shared with the
  logs and traces profiles (`namespace`, `app`), plus an `slo_service` selector
  label as used in the example (`http_requests_total{slo_service="example-service"}`).
- Generated rules live under `rules/slo/` in this repository for examples.
  Implementation repositories store their generated rules under their own
  `rules/slo/` path with matching `.promtool.yaml` mirrors.

## Phase 3 Entry Criteria

- Phase 0-2 Basic validation remains green.
- First implementation repository adoption is validated.
- Runtime secret and endpoint ownership remains outside this repository.
- Each profile has a rollback path and smoke test.

## Phase 3 Exit Criteria

- Each new profile has values, docs, validation, and fictional examples.
- `obsctl validate --strict-tools` covers added profile resources.
- Handover templates include profile-specific operational assumptions.
- Known limits are updated before implementation adoption.
