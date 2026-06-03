# Policy Samples

`pod-security-namespace.yaml` is referenced by the namespace Argo CD
Application. Other files are samples that implementation repositories should
review before applying.

Policy goals:

- keep the monitoring namespace explicit
- avoid unrestricted Argo CD destinations
- document hostPath exceptions
- keep notification and identity secrets outside this standard repository

## node-exporter and PodSecurity

`pod-security-namespace.yaml` sets `enforce: privileged` on the `monitoring`
namespace, and the `observability-namespaces` Argo CD Application applies it
before the stack syncs. This is deliberate: the kube-prometheus-stack
node-exporter DaemonSet needs host namespaces (`hostNetwork`, `hostPID`) and
hostPath volumes, which the `baseline` and `restricted` levels deny. Enforcing a
stricter level would block node-exporter at admission and drop node-level
metrics. `audit` and `warn` stay at `restricted`, so other workloads in the
namespace are still evaluated; node-exporter's own annotations are expected.

Implementations that isolate node-exporter (for example a dedicated namespace or
a cluster-level PodSecurity exemption for its service account) can tighten
`enforce` back to `baseline` for the rest of the monitoring workloads.
