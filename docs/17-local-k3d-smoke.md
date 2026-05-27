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

Expected evidence:

```powershell
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring get pods
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring get svc
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring get prometheusrules
helm --kubeconfig .tmp/kubeconfig/local-k3s.yaml -n monitoring status kube-prometheus-stack
```

## Delete Local k3s

```powershell
go run ./cmd/obsctl smoke local-k3s delete --name obs-standard-basic
```

## Notes

- This is still k3s, but packaged through k3d for local Docker-based testing.
- Do not use this as a production topology.
- If Docker Desktop is not running, start it before creating the cluster.
