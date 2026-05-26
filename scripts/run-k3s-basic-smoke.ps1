param(
  [Parameter(Mandatory = $true)]
  [string]$Kubeconfig,

  [string]$Namespace = "monitoring",
  [string]$ReleaseName = "kube-prometheus-stack"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $Kubeconfig)) {
  Write-Error "Kubeconfig not found: $Kubeconfig"
  exit 1
}

$helm = Get-Command helm -ErrorAction SilentlyContinue
if (-not $helm) {
  Write-Error "helm is required for the k3s smoke test."
  exit 1
}

$kubectl = Get-Command kubectl -ErrorAction SilentlyContinue
if (-not $kubectl) {
  Write-Error "kubectl is required for the k3s smoke test."
  exit 1
}

$env:KUBECONFIG = (Resolve-Path $Kubeconfig).Path

Write-Host "==> Validate repository"
$go = Get-Command go -ErrorAction SilentlyContinue
if ($go) {
  $goPath = $go.Source
} else {
  $portableGo = "C:\tmp\go-portable\go\bin\go.exe"
  if (Test-Path $portableGo) {
    $goPath = $portableGo
  }
}
if (-not $goPath) {
  Write-Error "go is required for strict repository validation."
  exit 1
}
& $goPath run ./cmd/obsctl validate --strict-tools

Write-Host "==> Cluster readiness"
kubectl get nodes

Write-Host "==> Create namespace"
kubectl create namespace $Namespace --dry-run=client -o yaml | kubectl apply -f -

Write-Host "==> Create placeholder Grafana admin Secret"
kubectl -n $Namespace create secret generic grafana-admin `
  --from-literal=admin-user=admin `
  --from-literal=admin-password=REPLACE_FOR_SMOKE_ONLY `
  --dry-run=client -o yaml | kubectl apply -f -

Write-Host "==> Helm dependency repo"
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

Write-Host "==> Install Basic profile"
helm upgrade --install $ReleaseName prometheus-community/kube-prometheus-stack `
  --version 85.0.2 `
  --namespace $Namespace `
  -f values/common/kube-prometheus-stack.yaml `
  -f values/profiles/basic.yaml `
  -f values/env/dev.yaml `
  -f values/sizing/small.yaml `
  --wait `
  --timeout 15m

Write-Host "==> Workload status"
kubectl -n $Namespace get pods
kubectl -n $Namespace get prometheusrules
kubectl -n $Namespace get servicemonitors

Write-Host "k3s Basic smoke test ok"
