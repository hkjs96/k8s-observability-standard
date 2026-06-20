# Phase 4 Plan

Phase 4 covers work that extends beyond the implemented Basic, Logs, Traces, and
SLO profiles. Everything here is optional: adopt an item only when an
implementation requirement justifies the added cost and operational complexity.
This document is the roadmap for the items listed as out of scope in
`.agent/context/product.md` and `docs/02-profiles.md`.

## Scope Principles

- Keep the implemented profiles (Basic, Logs, Traces, SLO) stable.
- Add one backend, generator, or UI at a time, each behind its own profile.
- Keep runtime secrets, endpoints, buckets, and identity bindings in
  implementation repositories.
- Add profile-scoped validation before publishing each new capability.
- Keep `charts-lock/chart-versions.yaml` as the single source of truth for any
  new pinned chart version.
- Prefer upstream chart defaults unless the standard needs an explicit guardrail.

## SLO Generation Pipeline (Sloth)

Goal: replace the hand-written generated-style PrometheusRule with a real Sloth
generation flow, as decided in `docs/01-standard-decisions.md` and
`docs/10-phase-3-plan.md`.

Planned outputs:

- Sloth spec schema guidance for the `examples/slo/` source specs.
- Documented `sloth generate` step that produces files under `rules/slo/`.
- Updated `.promtool.yaml` mirror generation for the produced rules.
- `obsctl validate profile slo` extension that checks the spec-to-rule mapping.

Open decisions:

- Whether Sloth runs in CI in implementation repositories or in this repository
  for examples only.
- Multi-window multi-burn-rate defaults versus the current page and warning
  pair.

## Pyrra SLO UI (Optional)

Goal: evaluate Pyrra as a dedicated SLO UI on top of the Sloth-generated rules.

Planned outputs:

- Pyrra values profile gated behind an explicit opt-in.
- Guidance on Pyrra reading existing generated rules versus owning generation.
- Validation that Pyrra does not duplicate or conflict with Sloth output.

Open decisions:

- Pyrra as read-only UI versus Pyrra as the generator (which would compete with
  Sloth).
- Whether a UI is justified for single-cluster implementations.

## Central / Long-Term Metrics (Mimir)

Goal: add central, long-term, or multi-cluster metrics storage when Prometheus
local retention is insufficient.

Planned outputs:

- Mimir values profile with object storage left as implementation-owned.
- Prometheus remote-write guidance pointing at an implementation-managed
  endpoint placeholder.
- Retention and tenancy guidance by sizing tier.

Open decisions:

- Single-tenant versus multi-tenant Mimir defaults.
- Where remote-write credentials and bucket configuration are owned.

## Dashboard-as-Code (Perses)

Goal: evaluate Perses for dashboard-as-code while keeping Grafana provisioning as
the initial UI.

Planned outputs:

- Perses PoC values profile.
- Conversion guidance for one existing dashboard candidate.
- Validation for committed dashboard-as-code definitions.

Open decisions:

- Perses alongside Grafana versus Perses as a replacement.
- How dashboard definitions stay in sync across both tools during a PoC.

## Profiling and eBPF (Optional, After Security Review)

Goal: evaluate continuous profiling and eBPF-based collectors after a security
review, since they need elevated privileges.

Planned outputs:

- Profiling backend candidate evaluation (for example Pyroscope).
- eBPF collector candidate evaluation (for example Beyla or Coroot).
- Pod Security, capability, and NetworkPolicy guardrails for privileged
  collectors.

Open decisions:

- Whether elevated privileges are acceptable in target clusters.
- Profiling backend storage and retention ownership.

## Phase 4 Entry Criteria

- Basic, Logs, Traces, and SLO validation remain green.
- The capability has a documented implementation requirement.
- Runtime secret, endpoint, bucket, and identity ownership remains outside this
  repository.
- Each capability has a rollback path and smoke test.

## Phase 4 Exit Criteria

- Each adopted capability has values, docs, validation, and fictional examples.
- `obsctl validate --strict-tools` covers added profile resources.
- Handover templates include capability-specific operational assumptions.
- Known limits and `.agent/context/product.md` scope are updated on adoption.
