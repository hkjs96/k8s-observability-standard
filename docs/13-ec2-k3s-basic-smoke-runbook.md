# EC2 k3s Basic Smoke Runbook

This runbook is a pre-production evidence guide for testing the Basic profile
on one disposable Amazon Linux 2023 EC2 instance with k3s. It is intentionally
generic. Keep all account, subnet, security group, key pair, host, and
kubeconfig values in the implementation repository or operator workstation.

## Scope

Use this runbook to confirm:

- Amazon Linux 2023 can bootstrap k3s with the example cloud-init.
- The Basic profile installs through Helm on a small single-node cluster.
- Prometheus, Alertmanager, Grafana, kube-state-metrics, and node-exporter pods
  become ready.
- The repository validation contract remains green before cluster install.

Do not use this runbook as a production topology. It is a disposable smoke test.

## Required Inputs

Prepare these values outside this repository:

| Input | Owner | Notes |
| --- | --- | --- |
| AWS region | Implementation repository | Use an approved test region. |
| AMI ID | Platform owner | Select Amazon Linux 2023 for the target region. |
| VPC and subnet | Platform owner | Prefer an isolated test subnet. |
| Security group | Platform owner | Restrict inbound access to operator source IPs. |
| SSH key pair | Operator | Do not commit key names or private keys. |
| Instance type | Operator | Use enough memory for kube-prometheus-stack. |
| Temporary kubeconfig path | Operator workstation | Store outside this repository when it contains a real host. |

## Security Group Minimums

Allow only the access needed for the test window:

| Direction | Port | Source | Purpose |
| --- | --- | --- | --- |
| Inbound | TCP 22 | Operator source IP | SSH and kubeconfig retrieval. |
| Inbound | TCP 6443 | Operator source IP | Kubernetes API, only if not using an SSH tunnel. |
| Outbound | TCP 443 | Any approved egress path | Package, k3s, Helm chart, and image pulls. |
| Outbound | TCP 80 | Any approved egress path | Only if the approved package path requires it. |

Prefer SSH tunneling or short-lived security group rules over broad public API
access. Do not expose Grafana or Prometheus directly for this smoke test.

## Provision

Launch one Amazon Linux 2023 instance with
`examples/ec2-k3s-basic/cloud-init.yaml` as user data.

The helper wraps `aws ec2 run-instances` and keeps all runtime values as
operator-provided parameters. It is exposed through `obsctl` so the same command
works on Windows, Linux, and macOS. Use `--dry-run` to inspect the AWS CLI call;
use `--yes` for a real launch:

```bash
go run ./cmd/obsctl smoke ec2-k3s launch \
  --region REPLACE_WITH_REGION \
  --ami-id REPLACE_WITH_AMI_ID \
  --instance-type REPLACE_WITH_INSTANCE_TYPE \
  --key-name REPLACE_WITH_KEY_PAIR_NAME \
  --subnet-id REPLACE_WITH_SUBNET_ID \
  --security-group-id REPLACE_WITH_SECURITY_GROUP_ID \
  --dry-run
```

After cloud-init completes, verify k3s on the instance:

```bash
sudo systemctl status k3s --no-pager
sudo kubectl get nodes -o wide
```

If the node is not ready, inspect:

```bash
sudo journalctl -u k3s --no-pager -n 200
sudo cloud-init status --long
```

## Kubeconfig

Copy kubeconfig from the instance:

```bash
sudo cat /etc/rancher/k3s/k3s.yaml
```

Or fetch it with the helper:

```bash
go run ./cmd/obsctl smoke ec2-k3s fetch-kubeconfig \
  --host REPLACE_WITH_INSTANCE_ADDRESS \
  --key-path REPLACE_WITH_KEY_PATH \
  --output REPLACE_WITH_KUBECONFIG_PATH
```

Store the copied file outside the repository if it contains a real endpoint.
The helper blocks repository-local output by default.
Replace the `server` value with either:

- `https://127.0.0.1:6443` when using an SSH tunnel.
- `https://REPLACE_WITH_TEMPORARY_INSTANCE_ADDRESS:6443` when the API port is
  temporarily restricted to the operator source IP.

Do not commit the rendered kubeconfig.

## Local Validation

From this repository, run strict validation before deployment:

```powershell
go run ./cmd/obsctl validate --strict-tools
```

If tools are not installed locally, install them with the shared helper first:

```bash
sh scripts/install-validation-tools.sh
```

On Windows runners or workstations:

```powershell
powershell -ExecutionPolicy Bypass -File scripts\install-validation-tools.ps1
```

## Smoke Install

Run the Basic smoke script against the temporary kubeconfig:

```bash
go run ./cmd/obsctl smoke k3s-basic install --kubeconfig REPLACE_WITH_KUBECONFIG_PATH
```

Expected cluster evidence:

```powershell
kubectl -n monitoring get pods
kubectl -n monitoring get svc
kubectl -n monitoring get prometheusrules
kubectl -n monitoring get servicemonitors
helm -n monitoring status kube-prometheus-stack
```

Record evidence in `templates/basic-smoke-evidence.template.md` from the
implementation repository copy. Do not paste real hostnames, account IDs, or
secret values into this standard repository.

## Cleanup

After evidence capture:

1. Delete or archive the temporary kubeconfig outside this repository.
2. Terminate the disposable EC2 instance.
3. Remove temporary security group ingress rules.
4. Record validation status, chart version, and known limits in the
   implementation repository.

The cleanup helper terminates the disposable instance. Use `--dry-run` to
inspect the AWS CLI call; use `--yes` for a real termination.

```bash
go run ./cmd/obsctl smoke ec2-k3s terminate \
  --region REPLACE_WITH_REGION \
  --instance-id REPLACE_WITH_INSTANCE_ID \
  --dry-run
```

## Failure Triage

Use this table to classify common failures:

| Symptom | First check | Likely owner |
| --- | --- | --- |
| EC2 reachable but node not ready | `sudo journalctl -u k3s --no-pager -n 200` | Platform or image baseline |
| Local kubectl cannot connect | Security group, SSH tunnel, kubeconfig `server` | Operator |
| Helm install times out | `kubectl -n monitoring get pods` | Standard values or node resources |
| Grafana pod fails | `kubectl -n monitoring describe pod -l app.kubernetes.io/name=grafana` | Values or resources |
| Prometheus rules rejected | `go run ./cmd/obsctl validate prometheus` | Standard rule validation |

Keep any implementation-specific remediation in the implementation repository.
