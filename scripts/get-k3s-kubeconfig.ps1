[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string]$HostName,

  [Parameter(Mandatory = $true)]
  [string]$KeyPath,

  [Parameter(Mandatory = $true)]
  [string]$OutputPath,

  [string]$User = "ec2-user",
  [string]$ServerHost,
  [int]$Port = 6443,
  [switch]$Force,
  [switch]$AllowRepositoryOutput
)

$ErrorActionPreference = "Stop"

$ssh = Get-Command ssh -ErrorAction SilentlyContinue
if (-not $ssh) {
  Write-Error "ssh is required."
  exit 1
}

if (-not (Test-Path $KeyPath)) {
  Write-Error "SSH key not found: $KeyPath"
  exit 1
}

if ((Test-Path $OutputPath) -and -not $Force) {
  Write-Error "OutputPath already exists. Use -Force to overwrite: $OutputPath"
  exit 1
}

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$resolvedOutputPath = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($OutputPath)
$normalizedRepoRoot = $repoRoot.TrimEnd("\", "/") + [System.IO.Path]::DirectorySeparatorChar
$isRepositoryOutput = $resolvedOutputPath.Equals($repoRoot, [System.StringComparison]::OrdinalIgnoreCase) -or
  $resolvedOutputPath.StartsWith($normalizedRepoRoot, [System.StringComparison]::OrdinalIgnoreCase)
if (-not $AllowRepositoryOutput -and $isRepositoryOutput) {
  Write-Error "OutputPath is inside this repository. Use a path outside the repository, or pass -AllowRepositoryOutput intentionally."
  exit 1
}

if (-not $ServerHost) {
  $ServerHost = $HostName
}

$outputDir = Split-Path -Parent $OutputPath
if ($outputDir -and -not (Test-Path $outputDir)) {
  New-Item -ItemType Directory -Path $outputDir | Out-Null
}

Write-Host "==> Fetch k3s kubeconfig"
$rawKubeconfig = & $ssh.Source -i $KeyPath -o StrictHostKeyChecking=accept-new "$User@$HostName" "sudo cat /etc/rancher/k3s/k3s.yaml"
if ($LASTEXITCODE -ne 0) {
  Write-Error "failed to fetch kubeconfig with ssh."
  exit $LASTEXITCODE
}

$server = "https://${ServerHost}:$Port"
$rendered = ($rawKubeconfig -join "`n") -replace "https://127\.0\.0\.1:6443", $server
$rendered = $rendered -replace "https://localhost:6443", $server

Set-Content -Path $OutputPath -Value $rendered -Encoding UTF8
Write-Host "Wrote kubeconfig: $OutputPath"
Write-Host "Do not commit this file if it contains a real endpoint."
