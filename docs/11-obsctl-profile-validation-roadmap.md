# obsctl Profile Validation Roadmap

This roadmap defines the intended validation contract before Phase 3 profile
resources are added.

## Current Contract

The current command validates the Basic baseline:

```powershell
go run ./cmd/obsctl validate --strict-tools
```

Current targets:

- `yaml`
- `basic`
- `argocd`
- `prometheus`
- `sensitive`
- `all`
- `profile basic`
- `profile logs`
- `profile traces`
- `profile slo`

All current profile targets are implemented for the resources present in this
repository.

## Future Contract

Phase 3 should add profile-specific validators without breaking the current
command:

```powershell
go run ./cmd/obsctl validate profile basic --strict-tools
go run ./cmd/obsctl validate profile logs --strict-tools
go run ./cmd/obsctl validate profile traces --strict-tools
go run ./cmd/obsctl validate profile slo --strict-tools
```

The existing `obsctl validate --strict-tools` command should continue to mean
"validate every profile that exists in this repository."

## Target Responsibilities

`profile basic`:

- YAML parse.
- Helm render and lint for `kube-prometheus-stack`.
- kubeconform rendered manifest validation.
- Prometheus rule mirror validation.
- Sensitive value scan.

`profile logs`:

- Loki values render.
- Alloy logs collector render.
- kubeconform rendered manifest validation.
- Static label-cardinality checks.
- Sensitive value scan.

`profile traces`:

- Tempo values render.
- OTLP receiver resource validation.
- OpenTelemetry Operator sample validation.
- Trace smoke test command documentation.
- Sensitive value scan.

`profile slo`:

- SLO spec linting.
- Generated PrometheusRule validation.
- `.promtool.yaml` mirror validation.
- Burn-rate rule promtool validation.
- Sensitive value scan.

## Implementation Notes

- Keep `validate.Options` as the shared options object.
- Add profile validators as separate functions instead of expanding `Basic`.
- Keep strict tool behavior consistent across profiles.
- Prefer structured profile metadata over path-only branching.
- Preserve the current single-target commands for compatibility.

## Exit Criteria

- Each profile has a target, tests, and documented required tools.
- `obsctl validate --strict-tools` runs all implemented profile targets.
- CI templates continue to call the all-profile strict command.
- Docs and `.agent` routing point to the new profile-specific targets when they exist.
