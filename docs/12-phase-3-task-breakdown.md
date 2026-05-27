# Phase 3 Task Breakdown

This document turns the Phase 3 plan into issue-ready work items. Keep each item
scoped to one profile and avoid implementation-owned runtime values.

## Logs Profile

### Task: Add Logs Profile Values

Goal: add reusable values for Loki and Alloy logs collection.

Deliverables:

- `values/profiles/logs.yaml`
- `values/profiles/logs-alloy.yaml`
- sizing guidance for logs storage and retention
- fictional example override

Acceptance criteria:

- Values render without implementation-owned secrets or endpoints.
- Labels follow the low-cardinality policy.
- `obsctl validate --strict-tools` remains green.

### Task: Add Logs Validation

Goal: validate logs profile resources before adoption.

Deliverables:

- `obsctl validate profile logs`
- static label-cardinality checks
- rendered manifest kubeconform check

Acceptance criteria:

- High-cardinality labels are rejected or explicitly documented.
- Missing runtime endpoints remain placeholders.
- Sensitive value scan remains green.

### Task: Add Logs Smoke Test

Goal: define implementation-owned log ingestion and query checks.

Deliverables:

- LogQL smoke query examples
- log ingestion checklist
- rollback notes

Acceptance criteria:

- Smoke tests do not require real endpoints in this repository.
- Implementation-owned values are clearly marked.

## Traces Profile

### Task: Add Traces Profile Values

Goal: add reusable values for Tempo and OTLP receiver workflows.

Deliverables:

- `values/profiles/traces.yaml`
- OTLP receiver resource guidance
- OpenTelemetry Operator sample

Acceptance criteria:

- Values avoid real endpoints and tenant identifiers.
- Trace components render with kubeconform where schemas are available.

### Task: Add Trace-To-Logs Guidance

Goal: document Grafana trace-to-logs expectations.

Deliverables:

- Grafana datasource guidance
- label mapping assumptions
- implementation-owned endpoint placeholders

Acceptance criteria:

- No real datasource URL is committed.
- Required labels are documented.

### Task: Add Traces Smoke Test

Goal: define implementation-owned trace ingestion checks.

Deliverables:

- OTLP test command example
- trace lookup checklist
- rollback notes

Acceptance criteria:

- Smoke test can be run from an implementation repository.
- This repository only contains placeholders.

## SLO Profile

### Task: Add SLO Spec Workflow

Goal: define SLO-as-code source and generated rule flow.

Deliverables:

- SLO spec location proposal
- generated PrometheusRule location proposal
- Sloth generation command documentation

Acceptance criteria:

- Generated rules have `.promtool.yaml` mirrors.
- Ownership of generated files is documented.

### Task: Add Burn-Rate Rule Validation

Goal: validate generated SLO alert rules.

Deliverables:

- promtool validation for generated rules
- SLO linting plan
- example availability and latency SLOs

Acceptance criteria:

- promtool passes for generated mirrors.
- Alert labels and annotations avoid implementation-owned values.

### Task: Add Error Budget Handover

Goal: make SLO operation handover repeatable.

Deliverables:

- error budget review template
- alert response notes
- rollback notes for generated rules

Acceptance criteria:

- Review template has owner, objective, budget window, and decision fields.
- No real service owner or escalation path is committed.

## Cross-Cutting Tasks

### Task: Add Profile-Scoped obsctl Targets

Goal: implement the validation roadmap without breaking current commands.

Deliverables:

- `obsctl validate profile basic`
- `obsctl validate profile logs`
- `obsctl validate profile traces`
- `obsctl validate profile slo`

Acceptance criteria:

- Existing `obsctl validate --strict-tools` still validates every implemented profile.
- Tests cover target parsing and profile dispatch.

### Task: Update Agent Routing

Goal: route future profile changes through profile-specific workflows and checks.

Deliverables:

- `.agent` workflow updates
- profile validation checks
- docs index updates

Acceptance criteria:

- Agents are directed to strict validation for profile changes.
- Basic profile guardrails remain intact.
