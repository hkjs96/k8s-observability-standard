# Kubernetes Observability Standard

This repository is the reusable standard baseline for Kubernetes observability
deployments. It is not a central multi-customer operations repository.

The first implementation covers Phase 0-2 of the v2 standard:

- standard decisions and profile documentation
- repository scaffold for Helm values, Argo CD, dashboards, rules, policies,
  templates, and examples
- Basic profile based on `kube-prometheus-stack`

Current status: Phase 0-2 Basic baseline is validation-ready. See
`docs/08-phase-0-2-completion.md` for completed outputs, validation evidence,
and known limits.

For a disposable Basic profile runtime check, use the Amazon Linux 2023 k3s
smoke-test material:

- `docs/13-ec2-k3s-basic-smoke-runbook.md`
- `examples/ec2-k3s-basic/`
- `templates/basic-smoke-evidence.template.md`

## Repository Model

Use this repository for common standards, sample values, validation commands,
and repeatable templates. Keep real deployment values in a separate
implementation repository.

Do not commit real account data, bucket names, domains, webhook URLs, IAM role
ARNs, credentials, or customer-specific cluster identifiers here.

## Layout

- `AGENTS.md`: Codex and generic agent entrypoint
- `CLAUDE.md`: Claude Code entrypoint
- `.agent/`: shared agent context, rules, workflows, and checks
- `docs/`: standard decisions, profiles, deployment, security, and operations
- `charts-lock/`: pinned upstream chart versions
- `values/`: common, profile, environment, and sizing values
- `argocd/`: AppProject and Application templates
- `dashboards/`: dashboard-as-code candidates and curated dashboard notes
- `rules/`: Prometheus recording and alerting rules
- `policies/`: security guardrail samples
- `templates/`: assessment, security review, deployment, and handover templates
- `examples/`: fictional example overlays and disposable smoke-test examples
- `cmd/obsctl`: provider-neutral validation CLI
- `scripts/validate-all.ps1`: fallback wrapper for environments without Go

## Agentic Workflow Model

Mandatory agent execution rules live in `.agent/`, not `docs/`.

- `.agent/index.md` routes changed paths to rules, workflows, and checks.
- `.agent/rules/` contains short mandatory constraints.
- `.agent/workflows/` contains task-specific execution steps.
- `.agent/checks/` maps validation criteria to scripts.
- `scripts/` contains executable validation used by agents, hooks, and CI.

`docs/` remains the human-facing standard and rationale.

## Phase 0-2 Validation

Run these checks before using a change as a deployable baseline:

```powershell
# List files.
rg --files

# Parse YAML if PyYAML is available.
python -c "import pathlib,yaml; [yaml.safe_load_all(p.read_text()) for p in pathlib.Path('.').rglob('*.yaml')]; print('yaml ok')"

# Preferred local check.
go run ./cmd/obsctl validate

# Fallback wrapper. Uses Go when available, otherwise runs legacy PowerShell validators.
powershell -ExecutionPolicy Bypass -File scripts/validate-all.ps1

# Render the Basic profile after installing Helm.
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm template kube-prometheus-stack prometheus-community/kube-prometheus-stack `
  --version 85.0.2 `
  --namespace monitoring `
  -f values/common/kube-prometheus-stack.yaml `
  -f values/profiles/basic.yaml `
  -f values/env/dev.yaml `
  -f values/sizing/small.yaml

# Validate rendered manifests after installing kubeconform.
helm template kube-prometheus-stack prometheus-community/kube-prometheus-stack `
  --version 85.0.2 `
  --namespace monitoring `
  -f values/common/kube-prometheus-stack.yaml `
  -f values/profiles/basic.yaml `
  -f values/env/dev.yaml `
  -f values/sizing/small.yaml |
  kubeconform -strict -ignore-missing-schemas

# Validate Prometheus rule mirrors after installing promtool.
promtool check rules rules/prometheus/*.promtool.yaml rules/alerting/*.promtool.yaml
```

## Sensitive Value Check

Use a manual review plus searches before publishing:

```powershell
rg -n --glob '!README.md' "adminPassword|access_key|secret_key|arn:aws|sourceRepos: \['\*'\]|namespace: '\*'"
rg -n "REPLACE_IN_IMPLEMENTATION_REPO"
```

Expected result: no committed real secrets or customer deployment values.
Placeholder values such as `REPLACE_IN_IMPLEMENTATION_REPO` must stay in
templates and must be replaced only in implementation repositories.
