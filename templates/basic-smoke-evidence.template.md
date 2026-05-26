# Basic Smoke Evidence

Use this template in an implementation repository after running a disposable
Basic profile smoke test. Leave real account, host, subnet, security group, and
secret values out of this standard repository.

## Test Context

- Standard repository revision:
- Chart version:
- Profile:
- Values layers:
- Kubernetes distribution:
- Operating system image family:
- Instance size class:
- Test start time:
- Test end time:
- Operator:

## Preflight Evidence

- `go run ./cmd/obsctl validate --strict-tools`:
- Sensitive value scan:
- Helm version:
- kubectl version:
- kubeconfig storage location owner:

## Cluster Evidence

- Node ready:
- `kubectl get nodes -o wide` captured:
- `kubectl -n monitoring get pods` captured:
- `kubectl -n monitoring get svc` captured:
- `kubectl -n monitoring get prometheusrules` captured:
- `kubectl -n monitoring get servicemonitors` captured:
- `helm -n monitoring status kube-prometheus-stack` captured:

## Functional Checks

- Prometheus pod ready:
- Alertmanager pod ready:
- Grafana pod ready:
- kube-state-metrics pod ready:
- node-exporter pod ready:
- Grafana datasource health checked:
- Prometheus targets checked:
- Dashboard discovery checked:

## Access And Cleanup

- SSH ingress limited to operator source:
- Kubernetes API ingress limited or tunneled:
- Grafana and Prometheus not publicly exposed:
- Temporary kubeconfig deleted or archived securely:
- EC2 instance terminated:
- Temporary security group ingress removed:

## Findings

- Passed:
- Failed:
- Known limits:
- Follow-up owner:
- Follow-up issue or ticket:
