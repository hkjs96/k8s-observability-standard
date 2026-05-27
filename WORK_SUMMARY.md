# Work Summary

## Done

- Updated local validation to use Go as the source of truth.
- Installed portable Go at `C:\tmp\go-portable\go\bin` and added it to the user PATH.
- Installed strict validation tools under `.tmp/tools`:
  - Helm `3.15.4`
  - kubeconform `0.6.7`
  - promtool `3.4.0`
- Removed legacy validation scripts from `scripts/`.
- Removed script-based validation references from README, docs, agent routing, and CI templates.
- Updated CI examples to install external tools inline and run `go run ./cmd/obsctl validate --strict-tools`.
- Fixed Windows smoke-test path handling for output paths on a different drive from the repository.
- Added profile-scoped validation targets:
  - `go run ./cmd/obsctl validate profile basic`
  - `go run ./cmd/obsctl validate profile logs`
- Added the Phase 3 logs profile baseline:
  - pinned Loki chart `7.0.0`
  - pinned Alloy chart `1.8.2`
  - reusable Loki values
  - reusable Alloy log collector values
  - fictional single-cluster logs override
  - logs profile documentation and smoke queries
- Added logs validation:
  - Helm render for Loki and Alloy
  - kubeconform validation for rendered manifests
  - static high-cardinality log label checks
- Updated `.agent` routing/checks to include profile-scoped validation.

## Validation

The following checks passed:

```powershell
go test ./...
go run ./cmd/obsctl validate --strict-tools
go run ./cmd/obsctl validate profile basic --strict-tools
go run ./cmd/obsctl validate profile logs --strict-tools
go run ./cmd/obsctl validate sensitive
```

## Next Work

1. Add traces profile values and validation.
2. Add SLO profile workflow and validation.
3. Decide whether logs profile should stay `SingleBinary` for small deployments or require implementation-managed object storage for production.
4. Add implementation repository adoption evidence once the first real adoption is tested.
5. Open a GitHub PR when GitHub CLI or connector-based PR creation is available.
