#!/usr/bin/env sh
set -eu

TOOLS_DIR=".tmp/tools"
mkdir -p "$TOOLS_DIR"

python3 -m pip install PyYAML --break-system-packages

HELM_VERSION="3.15.4"
curl -fsSL "https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz" -o /tmp/helm.tar.gz
tar -xzf /tmp/helm.tar.gz -C /tmp
mv /tmp/linux-amd64/helm "$TOOLS_DIR/helm"
chmod +x "$TOOLS_DIR/helm"

KUBECONFORM_VERSION="0.6.7"
curl -fsSL "https://github.com/yannh/kubeconform/releases/download/v${KUBECONFORM_VERSION}/kubeconform-linux-amd64.tar.gz" -o /tmp/kubeconform.tar.gz
tar -xzf /tmp/kubeconform.tar.gz -C "$TOOLS_DIR" kubeconform
chmod +x "$TOOLS_DIR/kubeconform"

PROMETHEUS_VERSION="3.4.0"
curl -fsSL "https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz" -o /tmp/prometheus.tar.gz
tar -xzf /tmp/prometheus.tar.gz -C /tmp
mv "/tmp/prometheus-${PROMETHEUS_VERSION}.linux-amd64/promtool" "$TOOLS_DIR/promtool"
chmod +x "$TOOLS_DIR/promtool"

echo "validation tools installed"
