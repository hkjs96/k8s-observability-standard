# Standard Decisions

## Phase 0 Decisions

| Area | Decision |
| --- | --- |
| First build | Start with `kube-prometheus-stack`. |
| Grafana | Use embedded Grafana from `kube-prometheus-stack` for Basic. |
| Dashboard management | Use ConfigMap labels, Grafana sidecar, and Argo CD. |
| Metrics backend | Use Prometheus by default. Add Mimir only when long-term, central, or multi-cluster metrics are required. |
| Logs backend | Keep Loki as the standard logs backend for later profiles. |
| Logs collector | Prefer Alloy for new standards. Fluent Bit is an implementation alternative when already mandated. |
| Promtail | Exclude from new standards; keep only migration notes outside Phase 0-2. |
| Traces backend | Keep Tempo for Advanced Traces, outside Phase 0-2. |
| SLO | Prefer Sloth for SLO-as-code proof of concept, outside Phase 0-2. |
| GitOps | Use upstream Helm charts, Git-hosted values layering, and Argo CD. |
| Security | Include namespace, RBAC, Pod Security, NetworkPolicy, Secret, and identity guardrails from the start. |
| Dashboard-as-code | Keep Grafana as the initial UI and evaluate Perses later as a PoC. |

## Non-Goals For Phase 0-2

- No full LGTM deployment in the first pass.
- No customer-specific deployment tree.
- No Mimir, Loki, Alloy, Tempo, Sloth, Pyrra, or Perses implementation yet.
- No real secret, endpoint, object storage, or identity values.

## Reuse Of Existing PoC Assets

Existing PoC assets can inform dashboard and rule selection, but they must be
reviewed before adoption. Do not directly copy:

- node-level service exposure
- embedded admin passwords
- real or legacy cluster labels
- object storage credentials
- PoC-only scripts or load-test manifests
- files with broken encoding
