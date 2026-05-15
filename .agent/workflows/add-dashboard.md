---
id: add-dashboard
type: workflow
required_for:
  - "dashboards/**"
checks:
  - go run ./cmd/obsctl validate sensitive
---

# Add Dashboard Workflow

1. Keep Grafana as the Phase 0-2 UI.
2. Sanitize dashboard JSON before committing.
3. Remove real datasource UIDs, cluster names, domains, account IDs, and customer labels.
4. Document dashboard purpose in `dashboards/grafana/README.md` when adding a new asset.
5. Run `go run ./cmd/obsctl validate sensitive`.
