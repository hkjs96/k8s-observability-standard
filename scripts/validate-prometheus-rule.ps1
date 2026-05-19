$ErrorActionPreference = "Stop"

& "$PSScriptRoot\validate-yaml.ps1" -Paths @("rules")

$ruleDirs = @("rules\prometheus", "rules\alerting")
foreach ($dir in $ruleDirs) {
  if (-not (Test-Path $dir)) {
    continue
  }
  $deployable = Get-ChildItem -Path $dir -Filter "*.yaml" -File |
    Where-Object { $_.Name -notlike "*.promtool.yaml" }
  foreach ($file in $deployable) {
    $mirror = Join-Path $file.DirectoryName ($file.BaseName + ".promtool.yaml")
    if (-not (Test-Path $mirror)) {
      Write-Error "Missing promtool mirror for $($file.FullName): $mirror"
      exit 1
    }
  }
}

$promtoolPath = $null
$localCandidates = @(
  ".tmp\tools\promtool.exe",
  "..\lgtm-k8s-observability-v2\tools\bin\promtool.exe"
)
foreach ($candidate in $localCandidates) {
  if (Test-Path $candidate) {
    $promtoolPath = (Resolve-Path $candidate).Path
    break
  }
}
if (-not $promtoolPath) {
  $expanded = Get-ChildItem -Path ".tmp\tools" -Recurse -Filter "promtool.exe" -File -ErrorAction SilentlyContinue |
    Select-Object -First 1
  if ($expanded) {
    $promtoolPath = $expanded.FullName
  }
}
if (-not $promtoolPath) {
  $promtool = Get-Command promtool -ErrorAction SilentlyContinue
  if ($promtool) {
    $promtoolPath = $promtool.Source
  }
}

if ($promtoolPath) {
  $mirrors = Get-ChildItem -Path "rules" -Recurse -Filter "*.promtool.yaml" -File
  if ($mirrors) {
    & $promtoolPath check rules @($mirrors.FullName)
    if ($LASTEXITCODE -ne 0) {
      exit $LASTEXITCODE
    }
  }
  Write-Host "promtool ok: $promtoolPath"
} else {
  Write-Host "promtool unavailable; skipped promtool check rules"
}

Write-Host "prometheus rules ok"
