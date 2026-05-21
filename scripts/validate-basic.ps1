$ErrorActionPreference = "Stop"

& "$PSScriptRoot\validate-yaml.ps1" -Paths @(".")

$helm = Get-Command helm -ErrorAction SilentlyContinue
if (-not $helm) {
  $localHelm = Join-Path (Get-Location) "..\lgtm-k8s-observability-v2\tools\bin\helm.exe"
  if (Test-Path $localHelm) {
    $helmPath = (Resolve-Path $localHelm).Path
  } else {
    Write-Host "helm unavailable; skipped helm template and helm lint"
    Write-Host "basic validation completed with helm skipped"
    exit 0
  }
} else {
  $helmPath = $helm.Source
}

$rendered = & $helmPath template kube-prometheus-stack kube-prometheus-stack `
  --repo https://prometheus-community.github.io/helm-charts `
  --version 85.0.2 `
  --namespace monitoring `
  -f values\common\kube-prometheus-stack.yaml `
  -f values\profiles\basic.yaml `
  -f values\env\dev.yaml `
  -f values\sizing\small.yaml

$kubeconformPath = $null
$kubeconformCandidates = @(
  ".tmp\tools\kubeconform.exe",
  "..\lgtm-k8s-observability-v2\tools\bin\kubeconform.exe"
)
foreach ($candidate in $kubeconformCandidates) {
  if (Test-Path $candidate) {
    $kubeconformPath = (Resolve-Path $candidate).Path
    break
  }
}
if (-not $kubeconformPath) {
  $kubeconform = Get-Command kubeconform -ErrorAction SilentlyContinue
  if ($kubeconform) {
    $kubeconformPath = $kubeconform.Source
  }
}
if ($kubeconformPath) {
  $rendered | & $kubeconformPath -strict -ignore-missing-schemas -summary
  if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
  }
} else {
  Write-Host "kubeconform unavailable; skipped rendered manifest validation"
}

$lintDir = Join-Path (Get-Location) (".cache\helm-lint-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Force -Path $lintDir | Out-Null

& $helmPath pull kube-prometheus-stack `
  --repo https://prometheus-community.github.io/helm-charts `
  --version 85.0.2 `
  --untar `
  --untardir $lintDir | Out-Null

& $helmPath lint (Join-Path $lintDir "kube-prometheus-stack") `
  -f values\common\kube-prometheus-stack.yaml `
  -f values\profiles\basic.yaml `
  -f values\env\dev.yaml `
  -f values\sizing\small.yaml

if ($LASTEXITCODE -ne 0) {
  exit $LASTEXITCODE
}

Write-Host "basic validation ok"
