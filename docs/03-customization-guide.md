# Customization Guide

Implementation repositories apply this standard by layering values files. This
repository provides reusable defaults and fictional examples only.

## Values Layering

Recommended order:

```yaml
helm:
  valueFiles:
    - $values/values/common/kube-prometheus-stack.yaml
    - $values/values/profiles/basic.yaml
    - $values/values/env/prod.yaml
    - $values/values/sizing/medium.yaml
    - $values/values/overrides/implementation-specific.yaml
```

The final implementation-specific override belongs in the implementation
repository, not this standard repository.

## Values Owned By Implementation Repos

- cluster identifier
- ingress hosts and TLS secret names
- object storage names and endpoints
- IAM role ARNs or workload identity bindings
- storage classes
- retention periods
- alert receivers and notification routes
- dashboard and rule exceptions
- runbook URLs

## Standard Defaults

Use the defaults here as a starting point:

- services are `ClusterIP` unless an implementation repo explicitly enables
  ingress or a private load balancer
- Grafana admin credentials come from an existing Secret
- Prometheus and Alertmanager run with conservative retention and storage
- dashboard ConfigMaps use `grafana_dashboard=1`
- alerting starts with a null receiver until real routing is provided elsewhere
