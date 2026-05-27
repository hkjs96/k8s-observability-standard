# SLO Profile

The SLO profile defines a reusable SLO-as-code workflow and generated
PrometheusRule validation. Real service ownership, escalation paths, and
notification routes remain implementation-owned.

## Artifacts

- `examples/slo/example-availability.slo.yaml`: fictional source SLO spec.
- `rules/slo/example-availability.yaml`: generated-style PrometheusRule sample.
- `rules/slo/example-availability.promtool.yaml`: promtool mirror.
- `templates/error-budget-review.template.md`: handover and review template.

## Workflow

1. Keep human-authored SLO specs in an implementation repository.
2. Generate PrometheusRule files from the SLO specs.
3. Commit generated rules with matching `.promtool.yaml` mirrors.
4. Run profile validation before merge.

```powershell
go run ./cmd/obsctl validate profile slo --strict-tools
go run ./cmd/obsctl validate --strict-tools
```

## Validation

The SLO validator checks required SLO artifacts, runs promtool through the
existing Prometheus rule validation path, and rejects implementation-owned
escalation values.

## Error Budget Review

Use `templates/error-budget-review.template.md` during handover and after major
incidents. Implementation repositories must fill in real owners, service names,
budget state, and decisions.

## Rollback

Rollback is implementation-owned. At minimum:

- Revert generated SLO rules if they page incorrectly.
- Keep SLO source and generated rules in sync.
- Review alert burn rates before re-enabling paging alerts.
