# Local k3d Smoke Test

This smoke test runs k3s locally through k3d and Docker. It avoids AWS, account
IDs, AMI IDs, subnets, security groups, and SSH keys.

## Requirements

- Docker Desktop running.
- `k3d` on `PATH`.
- `kubectl` on `PATH`.
- Helm available through this repository's `.tmp/tools/helm.exe` or on `PATH`.

## Create Local k3s

```powershell
go run ./cmd/obsctl smoke local-k3s create `
  --name obs-standard-basic `
  --kubeconfig .tmp/kubeconfig/local-k3s.yaml
```

The kubeconfig is written under `.tmp/`, which is ignored by git.
The helper rewrites k3d's `host.docker.internal` server URL to `127.0.0.1`
for Windows workstations where that hostname does not route back to the host.

## Install Basic Profile

```powershell
go run ./cmd/obsctl smoke k3s-basic install `
  --kubeconfig .tmp/kubeconfig/local-k3s.yaml
```

By default the smoke helper creates the Grafana admin Secret with:

- user: `admin`
- password: `REPLACE_FOR_SMOKE_ONLY`

You can override the password for a local run:

```powershell
go run ./cmd/obsctl smoke k3s-basic install `
  --kubeconfig .tmp/kubeconfig/local-k3s.yaml `
  --grafana-admin-password local-smoke-password
```

Expected evidence:

```powershell
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring get pods
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring get svc
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring get prometheusrules
helm --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring status kube-prometheus-stack
```

Optional Phase 0-2 artifact checks:

```powershell
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply -f dashboards/grafana/cluster-overview.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply -f dashboards/grafana/namespace-overview.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply -f rules/alerting/basic-alerts.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply -f rules/prometheus/basic-recording-rules.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply --dry-run=server -f policies/network-policy-sample.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply --dry-run=server -f policies/pod-security-namespace.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply --dry-run=server -f policies/rbac-readonly-sample.yaml
```

PrometheusRule samples include `release: kube-prometheus-stack` so the installed
Prometheus selector loads them during the smoke test.

## Delete Local k3s

```powershell
go run ./cmd/obsctl smoke local-k3s delete --name obs-standard-basic
```

## Notes

- This is still k3s, but packaged through k3d for local Docker-based testing.
- Do not use this as a production topology.
- If Docker Desktop is not running, start it before creating the cluster.
