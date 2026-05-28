# Phase 3 Smoke Workloads

These examples create disposable, fictional telemetry so the Phase 3 UI can be
checked in Grafana after Loki, Alloy, and Tempo are installed.

## Apply

```powershell
kubectl apply -f examples/phase3-smoke/namespace.yaml
kubectl apply -f examples/phase3-smoke/log-generator.yaml
kubectl apply -f examples/phase3-smoke/trace-generator.yaml
```

The same examples are applied by:

```powershell
go run ./cmd/obsctl smoke k3s-phase3 install --kubeconfig .tmp/kubeconfig/local-k3s.yaml
```

## Evidence

- `Logs Overview` should show log volume for namespace `observability-smoke`.
- `Logs Error-like Lines Detected` should have Loki data to evaluate.
- `Traces Overview` should show Tempo request activity and received spans.
- Grafana Explore should find traces for service `example-phase3-smoke`.

## Cleanup

```powershell
kubectl delete namespace observability-smoke
```
