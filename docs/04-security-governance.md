# Security Governance

Security review is part of the deployment design. It is not a follow-up cleanup
task.

## Required Guardrails

- Dedicated `monitoring` namespace or implementation-specific equivalent.
- Pod Security labels on the namespace.
- Least-privilege RBAC review before production.
- Argo CD AppProject with explicit source repositories and destinations.
- No `default` project usage for production.
- No wildcard source repository or wildcard destination namespace in production.
- Secrets are provided by Secret, ExternalSecret, SealedSecret, or SOPS in the
  implementation repository.
- No plaintext credentials in values files.

## Basic Profile Notes

`kube-prometheus-stack` deploys cluster-scoped resources and exporters. AppProject
permissions must account for Prometheus Operator CRDs, ClusterRoles, and
webhooks, while still avoiding broad wildcard access.

The sample namespace uses `privileged` enforcement with `restricted` audit/warn
labels. node-exporter requires host namespaces (`hostNetwork`/`hostPID`) and
hostPath volumes, which `baseline` and `restricted` deny, so a stricter enforced
level would block node-exporter at admission and drop node metrics. audit/warn
stay at `restricted` for visibility. See `policies/README.md` before changing the
enforced level (for example after isolating node-exporter).

## HostPath Policy

Phase 0-2 does not add a logs collector. Later logs profiles may allow a
read-only `/var/log` hostPath for the logs collector only.

Forbidden mounts for standard templates:

- `/`
- `/etc`
- `/var/run/docker.sock`
- `/run/containerd/containerd.sock`

Any hostPath exception must be documented in the security review template.
