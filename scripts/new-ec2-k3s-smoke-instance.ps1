[CmdletBinding(SupportsShouldProcess = $true, ConfirmImpact = "High")]
param(
  [Parameter(Mandatory = $true)]
  [string]$Region,

  [Parameter(Mandatory = $true)]
  [string]$AmiId,

  [Parameter(Mandatory = $true)]
  [string]$InstanceType,

  [Parameter(Mandatory = $true)]
  [string]$KeyName,

  [Parameter(Mandatory = $true)]
  [string]$SubnetId,

  [Parameter(Mandatory = $true)]
  [string]$SecurityGroupId,

  [string]$Name = "example-k3s-basic-smoke",
  [string]$Profile,
  [string]$UserDataPath = "examples\ec2-k3s-basic\cloud-init.yaml",
  [int]$VolumeSizeGiB = 40,
  [switch]$NoWait,
  [switch]$Force
)

$ErrorActionPreference = "Stop"

$aws = Get-Command aws -ErrorAction SilentlyContinue
if (-not $aws) {
  Write-Error "aws CLI is required."
  exit 1
}

if (-not (Test-Path $UserDataPath)) {
  Write-Error "User data file not found: $UserDataPath"
  exit 1
}

$resolvedUserData = (Resolve-Path $UserDataPath).Path
$blockDeviceMappings = @(
  @{
    DeviceName = "/dev/xvda"
    Ebs = @{
      VolumeSize = $VolumeSizeGiB
      VolumeType = "gp3"
      DeleteOnTermination = $true
    }
  }
) | ConvertTo-Json -Compress

$awsArgs = @(
  "ec2", "run-instances",
  "--region", $Region,
  "--image-id", $AmiId,
  "--instance-type", $InstanceType,
  "--key-name", $KeyName,
  "--subnet-id", $SubnetId,
  "--security-group-ids", $SecurityGroupId,
  "--user-data", "file://$resolvedUserData",
  "--metadata-options", "HttpTokens=required,HttpEndpoint=enabled",
  "--block-device-mappings", $blockDeviceMappings,
  "--tag-specifications", "ResourceType=instance,Tags=[{Key=Name,Value=$Name},{Key=Purpose,Value=k3s-basic-smoke}]",
  "--output", "json"
)

if ($Profile) {
  $awsArgs += @("--profile", $Profile)
}

$target = "$InstanceType from $AmiId in $SubnetId"
if (-not ($Force -or $PSCmdlet.ShouldProcess($target, "Launch EC2 k3s smoke instance"))) {
  Write-Host "Skipped launch."
  exit 0
}

Write-Host "==> Launch EC2 k3s smoke instance"
$runResult = & $aws.Source @awsArgs | ConvertFrom-Json
if ($LASTEXITCODE -ne 0) {
  Write-Error "aws ec2 run-instances failed."
  exit $LASTEXITCODE
}

$instanceId = $runResult.Instances[0].InstanceId
Write-Host "InstanceId: $instanceId"

if (-not $NoWait) {
  Write-Host "==> Wait for instance status ok"
  $waitArgs = @("ec2", "wait", "instance-status-ok", "--region", $Region, "--instance-ids", $instanceId)
  if ($Profile) {
    $waitArgs += @("--profile", $Profile)
  }
  & $aws.Source @waitArgs
  if ($LASTEXITCODE -ne 0) {
    Write-Error "aws ec2 wait instance-status-ok failed."
    exit $LASTEXITCODE
  }
}

$describeArgs = @(
  "ec2", "describe-instances",
  "--region", $Region,
  "--instance-ids", $instanceId,
  "--output", "json"
)
if ($Profile) {
  $describeArgs += @("--profile", $Profile)
}

$describeResult = & $aws.Source @describeArgs | ConvertFrom-Json
if ($LASTEXITCODE -ne 0) {
  Write-Error "aws ec2 describe-instances failed."
  exit $LASTEXITCODE
}

$instance = $describeResult.Reservations[0].Instances[0]
[pscustomobject]@{
  InstanceId = $instance.InstanceId
  State = $instance.State.Name
  PrivateIpAddress = $instance.PrivateIpAddress
  PublicIpAddress = $instance.PublicIpAddress
  Region = $Region
  Name = $Name
  KubeconfigCommand = "powershell -ExecutionPolicy Bypass -File scripts\get-k3s-kubeconfig.ps1 -HostName REPLACE_WITH_INSTANCE_ADDRESS -KeyPath REPLACE_WITH_KEY_PATH -OutputPath REPLACE_WITH_KUBECONFIG_PATH"
  CleanupCommand = "powershell -ExecutionPolicy Bypass -File scripts\remove-ec2-k3s-smoke-instance.ps1 -Region $Region -InstanceId $instanceId"
}
