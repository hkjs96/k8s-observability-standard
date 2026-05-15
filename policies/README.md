# Policy Samples

`pod-security-namespace.yaml` is referenced by the namespace Argo CD
Application. Other files are samples that implementation repositories should
review before applying.

Policy goals:

- keep the monitoring namespace explicit
- avoid unrestricted Argo CD destinations
- document hostPath exceptions
- keep notification and identity secrets outside this standard repository
