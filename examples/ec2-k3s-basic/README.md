# EC2 k3s Basic Smoke Example

This example is a disposable smoke-test pattern for validating the Basic
profile on a single Amazon Linux 2023 EC2 instance running k3s. It is
intentionally generic and does not include account-specific values.

## What This Tests

- k3s starts on one Amazon Linux 2023 EC2 instance.
- Helm can render and install `kube-prometheus-stack`.
- Grafana, Prometheus, Alertmanager, kube-state-metrics, and node-exporter pods start.
- Repository validation remains green before deployment.

## Required External Inputs

Provide these outside this repository:

- AWS region
- SSH key pair
- VPC and subnet
- security group rules
- EC2 instance type
- AMI ID
- temporary kubeconfig output path

Do not commit those values here.

## Suggested Flow

For the full evidence flow, use
`docs/13-ec2-k3s-basic-smoke-runbook.md`.

1. Launch one Amazon Linux 2023 EC2 instance with the cloud-init from
   `cloud-init.yaml`.
2. SSH to the instance and wait for k3s:

```bash
sudo systemctl status k3s --no-pager
sudo kubectl get nodes
```

Amazon Linux 2023 note: the cloud-init disables `nm-cloud-setup` when present
because it can interfere with Kubernetes CNI-managed routes on NetworkManager
based images.

3. Copy kubeconfig from the instance:

```bash
sudo cat /etc/rancher/k3s/k3s.yaml
```

4. Replace the server address in the copied kubeconfig with the instance private
   or temporary public address.
5. From this repository, run validation:

```powershell
go run ./cmd/obsctl validate --strict-tools
```

6. Install the Basic profile using the smoke script:

```powershell
powershell -ExecutionPolicy Bypass -File scripts\run-k3s-basic-smoke.ps1 -Kubeconfig REPLACE_WITH_KUBECONFIG_PATH
```

7. Destroy the EC2 instance when testing is complete.

## Expected Result

- `monitoring` namespace exists.
- `kube-prometheus-stack` Helm release is deployed.
- Prometheus Operator pods are ready.
- Grafana pod is ready.
- ServiceMonitors and PrometheusRules are accepted by the cluster.

## Cleanup

Use the implementation-owned infrastructure tool to terminate the EC2 instance.
This repository does not own the EC2 lifecycle.

Also remove temporary security group ingress rules and delete any kubeconfig
copy that contains a real instance address.
