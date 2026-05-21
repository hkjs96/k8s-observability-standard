---
id: phase-3-profile-validation
type: check
applies_to:
  - "docs/10-phase-3-plan.md"
  - "values/**"
  - "rules/**"
  - "dashboards/**"
  - "policies/**"
commands:
  - "go run ./cmd/obsctl validate --strict-tools"
  - "go run ./cmd/obsctl validate sensitive"
---

# Phase 3 Profile Validation

Run:

```powershell
go run ./cmd/obsctl validate --strict-tools
go run ./cmd/obsctl validate sensitive
```

Completion criteria:

- Basic profile validation remains green.
- Any new profile resources render or have a documented validation gap.
- Prometheus rules include matching `.promtool.yaml` mirrors.
- No runtime-specific secret, endpoint, account, bucket, or identity value is committed.
- Phase 3 changes remain scoped to logs, tracing, SLO, or explicitly documented optional profiles.
