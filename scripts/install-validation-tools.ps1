$ErrorActionPreference = "Stop"

$tools = ".tmp\tools"
New-Item -ItemType Directory -Force -Path $tools | Out-Null

python -m pip install PyYAML

$kubeconformVersion = "0.6.7"
$kubeconformZip = Join-Path $tools "kubeconform.zip"
Invoke-WebRequest `
  -Uri "https://github.com/yannh/kubeconform/releases/download/v$kubeconformVersion/kubeconform-windows-amd64.zip" `
  -OutFile $kubeconformZip
Expand-Archive -Path $kubeconformZip -DestinationPath $tools -Force

$prometheusVersion = "3.4.0"
$prometheusZip = Join-Path $tools "prometheus.zip"
Invoke-WebRequest `
  -Uri "https://github.com/prometheus/prometheus/releases/download/v$prometheusVersion/prometheus-$prometheusVersion.windows-amd64.zip" `
  -OutFile $prometheusZip
Expand-Archive -Path $prometheusZip -DestinationPath $tools -Force

Write-Host "validation tools installed"
