# Single Cluster Basic Example

This example shows how an implementation repository could apply the Basic
profile for one cluster. Values are fictional and must be replaced outside this
standard repository.

Recommended value file order:

```yaml
helm:
  valueFiles:
    - $values/values/common/kube-prometheus-stack.yaml
    - $values/values/profiles/basic.yaml
    - $values/values/env/dev.yaml
    - $values/values/sizing/small.yaml
    - $values/values/overrides/single-cluster-basic.yaml
```

Required implementation-owned resources:

- `grafana-admin` Secret in `monitoring`
- alert receiver Secret if notifications are enabled
- ingress or private access configuration if Grafana is exposed
