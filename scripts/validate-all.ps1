$ErrorActionPreference = "Stop"

$go = Get-Command go -ErrorAction SilentlyContinue
if ($go) {
  $goPath = $go.Source
} else {
  $portableGo = "C:\tmp\go-portable\go\bin\go.exe"
  if (Test-Path $portableGo) {
    $goPath = $portableGo
  }
}

if ($goPath) {
  Write-Host "==> go run ./cmd/obsctl validate"
  $cacheRoot = Join-Path (Resolve-Path ".cache").Path "go-build"
  $modRoot = Join-Path (Resolve-Path ".cache").Path "go-mod"
  $env:GOCACHE = $cacheRoot
  $env:GOMODCACHE = $modRoot
  & $goPath run ./cmd/obsctl validate
  exit $LASTEXITCODE
}

Write-Host "go unavailable; falling back to legacy PowerShell validators"

$checks = @(
  "validate-basic.ps1",
  "validate-argocd.ps1",
  "validate-prometheus-rule.ps1",
  "validate-sensitive-values.ps1"
)

foreach ($check in $checks) {
  $path = Join-Path $PSScriptRoot $check
  Write-Host "==> $check"
  & $path
}

Write-Host "all validation ok"
