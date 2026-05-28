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
- Added the Phase 3 traces profile baseline:
  - pinned Tempo chart `1.24.4`
  - reusable Tempo values
  - fictional single-cluster traces override
  - OpenTelemetry Operator instrumentation sample
  - traces profile documentation and smoke checks
- Added traces validation:
  - Helm render for Tempo
  - kubeconform validation for rendered manifests
  - static implementation-owned endpoint and tenant checks
- Added the Phase 3 SLO profile workflow:
  - fictional source SLO spec
  - generated-style PrometheusRule sample
  - promtool mirror
  - error budget review template
- Added SLO validation:
  - required artifact checks
  - promtool validation through the existing rule validator
  - static implementation-owned escalation checks
- Updated `.agent` routing/checks to include profile-scoped validation.
- Added local Docker-based k3s smoke commands through k3d:
  - `go run ./cmd/obsctl smoke local-k3s create`
  - `go run ./cmd/obsctl smoke local-k3s delete`
  - `docs/17-local-k3d-smoke.md`
- Ran local k3s smoke on this workstation:
  - Docker Desktop started successfully.
  - k3d cluster `obs-standard-basic` created.
  - k3s node reached `Ready`.
  - Basic `kube-prometheus-stack` Helm release deployed in `monitoring`.
  - Grafana smoke login is created by the helper as `admin` / `REPLACE_FOR_SMOKE_ONLY` unless overridden with `--grafana-admin-password`.
  - Grafana, Prometheus, Alertmanager, kube-state-metrics, and node-exporter pods reached `Running`.
  - Prometheus, Grafana, and Alertmanager HTTP readiness checks passed.
  - Prometheus target API reported all active targets up.
  - Grafana API found the standard Cluster Overview and Namespace Overview dashboards after applying dashboard ConfigMaps.
  - Added Grafana-managed alert provisioning for logs/traces profile checks.
  - Added Loki and Tempo datasource provisioning plus Logs Overview and Traces Overview dashboards.
  - Installed Loki, Alloy, and Tempo into local k3s and confirmed Grafana sees Loki/Tempo datasources, Logs/Traces dashboards, and two Grafana-managed alert rules.
  - Added Phase 3 smoke workloads for generated logs, generated traces, and example SLO metrics.
  - Added an SLO Overview dashboard for the example availability SLO.
  - Ran the Phase 3 local smoke workload and confirmed Loki query success, Tempo received/accepted spans at 484, and the example SLO ratio at 0.97.
  - PrometheusRule samples were applied and updated with the release label required by the installed Prometheus selector.
  - Policy samples passed server-side dry-run validation; PodSecurity baseline warns on node-exporter host namespace/hostPath/hostPort usage, which is expected for the node-exporter daemonset.

## Validation

The following checks passed:

```powershell
go test ./...
go run ./cmd/obsctl validate --strict-tools
go run ./cmd/obsctl validate profile basic --strict-tools
go run ./cmd/obsctl validate profile logs --strict-tools
go run ./cmd/obsctl validate profile traces --strict-tools
go run ./cmd/obsctl validate profile slo --strict-tools
go run ./cmd/obsctl validate sensitive
```

## Next Work

1. Decide whether logs profile should stay `SingleBinary` for small deployments or require implementation-managed object storage for production.
2. Decide whether traces profile should keep the deprecated single binary Tempo chart or move to another deployment model before production adoption.
3. Decide whether SLO generation should standardize on Sloth, Pyrra, or provider-specific tooling.
4. Add implementation repository adoption evidence once the first real adoption is tested.
