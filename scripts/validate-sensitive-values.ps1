$ErrorActionPreference = "Stop"

$rg = Get-Command rg -ErrorAction SilentlyContinue
if (-not $rg) {
  Write-Error "rg is required for sensitive value validation."
  exit 1
}

$patterns = @(
  "adminPassword",
  "access_key",
  "secret_key",
  "arn:aws",
  "sourceRepos: \['\*'\]",
  "namespace: '\*'",
  "NodePort",
  "admin123",
  "shi-cluster",
  "loki-chunks",
  "ACCESS_KEY",
  "SECRET_KEY"
)

$combined = $patterns -join "|"
$output = & rg -n --glob "!.cache/**" --glob "!README.md" --glob "!scripts/validate-sensitive-values.ps1" --glob "!internal/validate/sensitive.go" $combined . 2>$null
if ($LASTEXITCODE -eq 0) {
  Write-Error "Forbidden sensitive or implementation-specific pattern found:`n$output"
  exit 1
}
if ($LASTEXITCODE -gt 1) {
  Write-Error "rg failed while scanning sensitive values."
  exit 1
}

$customerMentions = & rg -n --glob "!.cache/**" --glob "!scripts/validate-sensitive-values.ps1" --glob "!internal/validate/sensitive.go" "customers/" . 2>$null
if ($LASTEXITCODE -eq 0) {
  $bad = $customerMentions | Where-Object {
    ($_ -notmatch "AGENTS\.md") -and
    ($_ -notmatch "repository-boundary\.md") -and
    ($_ -notmatch "00-overview\.md")
  }
  if ($bad) {
    Write-Error "Unexpected customers/ mention found:`n$($bad -join [Environment]::NewLine)"
    exit 1
  }
} elseif ($LASTEXITCODE -gt 1) {
  Write-Error "rg failed while scanning customers/ mentions."
  exit 1
}

Write-Host "sensitive values ok"
