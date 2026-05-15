$ErrorActionPreference = "Stop"

& "$PSScriptRoot\validate-yaml.ps1" -Paths @("argocd")

$rg = Get-Command rg -ErrorAction SilentlyContinue
if (-not $rg) {
  Write-Error "rg is required for Argo CD validation."
  exit 1
}

$defaultProject = & rg -n "project:\s+default\b" argocd 2>$null
if ($LASTEXITCODE -eq 0) {
  Write-Error "Argo CD Application must not use the default project:`n$defaultProject"
  exit 1
}
if ($LASTEXITCODE -gt 1) {
  Write-Error "rg failed while checking Argo CD projects."
  exit 1
}

$wildcard = & rg -n "sourceRepos:\s*\['\*'\]|namespace:\s*'\*'" argocd 2>$null
if ($LASTEXITCODE -eq 0) {
  Write-Error "Argo CD wildcard sourceRepos or destination namespace found:`n$wildcard"
  exit 1
}
if ($LASTEXITCODE -gt 1) {
  Write-Error "rg failed while checking Argo CD wildcards."
  exit 1
}

$apps = Get-ChildItem -Path "argocd\applications" -Filter "*.yaml" -File -ErrorAction SilentlyContinue
foreach ($app in $apps) {
  $content = Get-Content -Raw -Path $app.FullName
  if ($content -match "chart:\s+" -and $content -notmatch "valueFiles:") {
    Write-Error "Helm Application missing valueFiles: $($app.FullName)"
    exit 1
  }
}

Write-Host "argocd ok"
