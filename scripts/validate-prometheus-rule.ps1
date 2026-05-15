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

$promtool = Get-Command promtool -ErrorAction SilentlyContinue
if ($promtool) {
  $mirrors = Get-ChildItem -Path "rules" -Recurse -Filter "*.promtool.yaml" -File
  if ($mirrors) {
    & $promtool.Source check rules @($mirrors.FullName)
    if ($LASTEXITCODE -ne 0) {
      exit $LASTEXITCODE
    }
  }
  Write-Host "promtool ok"
} else {
  Write-Host "promtool unavailable; skipped promtool check rules"
}

Write-Host "prometheus rules ok"
