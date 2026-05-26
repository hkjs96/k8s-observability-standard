# Day2 Operations

Day2 operations are owned by the implementation repository, but this standard
defines the baseline checks.

## Recurring Checks

- Prometheus TSDB head series and sample ingestion rate.
- Target health and scrape duration.
- Prometheus and Alertmanager PVC usage if persistence is enabled.
- Grafana datasource health.
- Rule evaluation failures.
- Alert route test after notification receiver changes.
- AppProject drift and unexpected permissions.

## Change Control

Review before applying:

- chart version upgrades
- CRD upgrades
- Prometheus retention changes
- persistent volume resizing
- rule changes that affect paging behavior
- dashboard provisioning format changes

## Incident Handover Minimum

Each implementation repo should maintain:

- owner and escalation path
- dashboard entry points
- alert route map
- smoke test commands
- latest smoke-test evidence location
- rollback procedure for values changes
- known limits and sizing assumptions
