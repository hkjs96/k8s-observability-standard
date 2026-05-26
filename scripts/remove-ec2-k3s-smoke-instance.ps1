[CmdletBinding(SupportsShouldProcess = $true, ConfirmImpact = "High")]
param(
  [Parameter(Mandatory = $true)]
  [string]$Region,

  [Parameter(Mandatory = $true)]
  [string]$InstanceId,

  [string]$Profile,
  [switch]$NoWait,
  [switch]$Force
)

$ErrorActionPreference = "Stop"

$aws = Get-Command aws -ErrorAction SilentlyContinue
if (-not $aws) {
  Write-Error "aws CLI is required."
  exit 1
}

$target = "$InstanceId in $Region"
if (-not ($Force -or $PSCmdlet.ShouldProcess($target, "Terminate EC2 instance"))) {
  Write-Host "Skipped termination."
  exit 0
}

$terminateArgs = @(
  "ec2", "terminate-instances",
  "--region", $Region,
  "--instance-ids", $InstanceId,
  "--output", "json"
)
if ($Profile) {
  $terminateArgs += @("--profile", $Profile)
}

Write-Host "==> Terminate EC2 k3s smoke instance"
& $aws.Source @terminateArgs
if ($LASTEXITCODE -ne 0) {
  Write-Error "aws ec2 terminate-instances failed."
  exit $LASTEXITCODE
}

if (-not $NoWait) {
  Write-Host "==> Wait for instance terminated"
  $waitArgs = @("ec2", "wait", "instance-terminated", "--region", $Region, "--instance-ids", $InstanceId)
  if ($Profile) {
    $waitArgs += @("--profile", $Profile)
  }
  & $aws.Source @waitArgs
  if ($LASTEXITCODE -ne 0) {
    Write-Error "aws ec2 wait instance-terminated failed."
    exit $LASTEXITCODE
  }
}

Write-Host "EC2 k3s smoke instance cleanup ok"
