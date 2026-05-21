# Phase 3 Plan

Phase 3 expands beyond the Basic metrics baseline. The work should remain split
by profile so logs, tracing, and SLO features can be adopted independently.

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
- Storage and retention guidance.

Validation additions:

- Helm render and kubeconform for log components.
- Static check for forbidden high-cardinality labels.
- Example LogQL smoke queries.

Open decisions:

- Single-cluster Loki versus implementation-managed central backend.
- Object storage requirement boundary.
- Default retention guidance by sizing tier.

## Advanced Traces Profile

Goal: add distributed tracing and trace-to-logs workflows.

Planned outputs:

- Tempo values profile.
- OTLP receiver Service and routing guidance.
- OpenTelemetry Operator instrumentation samples.
- Grafana trace-to-logs settings.
- Trace smoke test.

Validation additions:

- Helm render and kubeconform for trace components.
- Static checks for OTLP endpoint placeholders.
- Example trace ingestion smoke command.

Open decisions:

- Tempo storage mode for small deployments.
- Required versus optional OpenTelemetry Operator.
- Default sampling guidance.

## SLO Profile

Goal: define SLO-as-code workflows and error budget operations.

Planned outputs:

- Sloth rule generation workflow.
- Availability and latency SLO examples.
- Alert routing guidance for burn-rate alerts.
- Error budget review template.

Validation additions:

- SLO spec linting.
- Generated PrometheusRule mirror validation.
- promtool check for generated burn-rate rules.

Open decisions:

- Sloth only versus Pyrra UI option.
- Standard service labels required for SLO selection.
- Where generated rules are stored in implementation repositories.

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
