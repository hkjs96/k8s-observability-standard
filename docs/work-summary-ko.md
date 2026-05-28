# 작업 및 검증 결과 요약

## 작업 범위

이 문서는 현재 저장소에서 실제로 처리한 작업과 검증 결과를 한국어로 정리한 기록입니다.

대상 저장소는 Kubernetes Observability Standard이며, Phase 0-2 Basic 기준 검증과 Phase 3 profile 확장을 함께 진행했습니다.

## 처리한 내용

### Go 기반 검증 체계 정리

- 기존 Windows/PowerShell 기반 검증 스크립트를 제거했습니다.
- `cmd/obsctl` Go CLI를 검증의 기준으로 정리했습니다.
- CI 예시는 외부 도구 설치 후 `go run ./cmd/obsctl validate --strict-tools`를 실행하도록 수정했습니다.
- portable Go를 설치하고 로컬 검증에 사용했습니다.
  - 경로: `C:\tmp\go-portable\go\bin`
  - 버전: `go1.26.3 windows/amd64`

### 검증 도구 설치

`.tmp/tools` 아래에 strict validation용 도구를 설치했습니다.

- Helm `3.15.4`
- kubeconform `0.6.7`
- promtool `3.4.0`
- k3d `5.7.5`

### Profile 검증 대상 추가

다음 profile-scoped 검증 명령을 추가했습니다.

```powershell
go run ./cmd/obsctl validate profile basic
go run ./cmd/obsctl validate profile logs
go run ./cmd/obsctl validate profile traces
go run ./cmd/obsctl validate profile slo
```

전체 검증 명령인 아래 명령에도 logs, traces, slo 검증이 포함되도록 확장했습니다.

```powershell
go run ./cmd/obsctl validate --strict-tools
```

### Logs Profile 추가

- Loki chart pin 추가
  - chart: `grafana/loki`
  - version: `7.0.0`
- Alloy chart pin 추가
  - chart: `grafana/alloy`
  - version: `1.8.2`
- 추가 파일:
  - `values/profiles/logs.yaml`
  - `values/profiles/logs-alloy.yaml`
  - `values/overrides/single-cluster-logs.yaml`
  - `docs/14-logs-profile.md`
- 검증 내용:
  - Loki Helm render
  - Alloy Helm render
  - kubeconform manifest 검증
  - 고카디널리티 log label 정적 검사

### Traces Profile 추가

- Tempo chart pin 추가
  - chart: `grafana/tempo`
  - version: `1.24.4`
- 추가 파일:
  - `values/profiles/traces.yaml`
  - `values/overrides/single-cluster-traces.yaml`
  - `examples/opentelemetry/traces-instrumentation.yaml`
  - `docs/15-traces-profile.md`
- 검증 내용:
  - Tempo Helm render
  - kubeconform manifest 검증
  - 실제 vendor endpoint, tenant 값 등 implementation-owned 값 정적 검사
- 참고:
  - 사용한 Tempo single binary chart는 upstream에서 deprecated 상태라 production adoption 전에 재검토가 필요합니다.

### SLO Profile 추가

- 추가 파일:
  - `examples/slo/example-availability.slo.yaml`
  - `rules/slo/example-availability.yaml`
  - `rules/slo/example-availability.promtool.yaml`
  - `templates/error-budget-review.template.md`
  - `docs/16-slo-profile.md`
- 검증 내용:
  - 필수 SLO artifact 존재 확인
  - promtool rule 검증
  - 실제 escalation 값, webhook, 계정 식별자 등 implementation-owned 값 정적 검사

### Local k3s Smoke Test 추가

AWS 없이 로컬 Docker 기반으로 k3s를 테스트할 수 있도록 k3d smoke path를 추가했습니다.

추가 명령:

```powershell
go run ./cmd/obsctl smoke local-k3s create
go run ./cmd/obsctl smoke local-k3s delete
```

추가 문서:

- `docs/17-local-k3d-smoke.md`

추가 개선:

- k3d kubeconfig의 `host.docker.internal` 주소를 Windows 환경에서 접근 가능한 `127.0.0.1`로 rewrite하도록 처리했습니다.
- Grafana smoke 계정 정보를 명확히 문서화했습니다.
- Grafana smoke password를 `--grafana-admin-password`로 override할 수 있게 했습니다.

기본 Grafana smoke login:

- URL: `http://127.0.0.1:3000`
- user: `admin`
- password: `REPLACE_FOR_SMOKE_ONLY`

## 실제 수행한 테스트

### Go Unit Test

```powershell
go test ./...
```

결과:

- 통과
- 대상 package:
  - `cmd/obsctl`
  - `internal/smoke`
  - `internal/validate`
  - `internal/walk`

### Strict Validation

```powershell
go run ./cmd/obsctl validate --strict-tools
```

결과:

- 통과

확인된 항목:

- YAML parse 통과
- Basic Helm render 통과
- Basic Helm lint 통과
- kubeconform manifest 검증 통과
- Logs profile render 및 kubeconform 통과
- Traces profile render 및 kubeconform 통과
- SLO promtool 검증 통과
- Argo CD template 검증 통과
- PrometheusRule promtool 검증 통과
- sensitive value scan 통과

### Profile 개별 검증

```powershell
go run ./cmd/obsctl validate profile basic --strict-tools
go run ./cmd/obsctl validate profile logs --strict-tools
go run ./cmd/obsctl validate profile traces --strict-tools
go run ./cmd/obsctl validate profile slo --strict-tools
go run ./cmd/obsctl validate sensitive
```

결과:

- 모두 통과

## 실제 Local k3s 구축 결과

Docker Desktop과 k3d를 사용해 로컬 k3s cluster를 구축했습니다.

생성한 cluster:

- 이름: `obs-standard-basic`
- k3s version: `v1.29.6+k3s1`
- node: `k3d-obs-standard-basic-server-0`
- 상태: `Ready`

Basic profile 설치:

- Helm release: `kube-prometheus-stack`
- namespace: `monitoring`
- 상태: `deployed`

실행한 설치 명령:

```powershell
go run ./cmd/obsctl smoke k3s-basic install --kubeconfig .tmp/kubeconfig/local-k3s.yaml --skip-validation
```

## Kubernetes 리소스 확인 결과

`monitoring` namespace에서 다음 pod들이 Running/Ready 상태였습니다.

- `alertmanager-kube-prometheus-stack-alertmanager-0`
- `kube-prometheus-stack-grafana`
- `kube-prometheus-stack-kube-state-metrics`
- `kube-prometheus-stack-operator`
- `kube-prometheus-stack-prometheus-node-exporter`
- `prometheus-kube-prometheus-stack-prometheus-0`

확인 결과:

- Deployment ready
- StatefulSet ready
- DaemonSet ready
- Prometheus CR `Available=True`
- Alertmanager CR `Available=True`

## Grafana 확인 결과

Grafana port-forward를 열고 실제 API health와 dashboard 검색을 확인했습니다.

접속 정보:

- URL: `http://127.0.0.1:3000`
- user: `admin`
- password: `REPLACE_FOR_SMOKE_ONLY`

확인 결과:

- Grafana database: `ok`
- Grafana version: `13.0.1`
- dashboard 검색 결과 확인:
  - `Cluster Overview`
  - `Namespace Overview`
  - `Prometheus / Overview`
  - `Alertmanager / Overview`
  - `Grafana Overview`

추가로 logs/traces profile UI를 위해 Grafana provisioning 구조를 적용했습니다.

운영 기준:

- Kubernetes/platform 핵심 알람은 `PrometheusRule` + Alertmanager로 관리합니다.
- Dashboard 기반 logs/traces/application-facing 알람은 Grafana-managed alert로 관리합니다.
- Grafana 알람도 UI에서만 만들지 않고 ConfigMap provisioning으로 Git 관리합니다.

실제 local k3s에 Loki, Alloy, Tempo를 설치한 뒤 Grafana API로 확인한 결과:

- datasource:
  - `Prometheus`
  - `Alertmanager`
  - `Loki`
  - `Tempo`
- dashboard:
  - `Cluster Overview`
  - `Namespace Overview`
  - `Logs Overview`
  - `Traces Overview`
  - `Prometheus / Overview`
  - `Alertmanager / Overview`
- Grafana-managed alert rule:
  - `Logs Error-like Lines Detected`
  - `Tempo Targets Down`

즉, logs/traces는 PrometheusRule만으로 억지로 처리하지 않고 Grafana UI에서 직접 확인하고 운영할 수 있게 구성했습니다.

추가로 Phase 3 화면에 실제 샘플 값이 나오도록 local smoke workload를 추가했습니다.

- `examples/phase3-smoke/log-generator.yaml`
  - Loki/Logs Overview 확인용 로그를 지속 출력합니다.
  - 주기적으로 `error` 문자열을 출력해서 Grafana log alert 평가에도 데이터가 생깁니다.
- `examples/phase3-smoke/trace-generator.yaml`
  - OpenTelemetry `telemetrygen`으로 Tempo에 trace span을 전송합니다.
  - 이 workload가 실행된 뒤 `Traces Overview`의 request/spans 패널에 값이 생깁니다.
  - local k3s에서 확인한 값:
    - `sum(tempo_distributor_spans_received_total)` = `484`
    - `sum(tempo_receiver_accepted_spans)` = `484`
- `examples/phase3-smoke/slo-metrics-generator.yaml`
  - `http_requests_total{slo_service="example-service"}` 샘플 metric을 노출합니다.
  - 기존 SLO recording rule과 새 `SLO Overview` dashboard 확인에 사용합니다.
- `dashboards/grafana/slo-overview.yaml`
  - example availability, request rate, error budget burn을 확인하는 Grafana dashboard입니다.
  - local k3s에서 확인한 값:
    - `slo:example_service_availability:ratio_rate5m` = `0.9700000000000001`
    - burn rate query 결과 = `29.99999999999989`

한 번에 적용하는 명령:

```powershell
go run ./cmd/obsctl smoke k3s-phase3 install --kubeconfig .tmp/kubeconfig/local-k3s.yaml
```

## Prometheus 확인 결과

Prometheus API를 port-forward로 확인했습니다.

확인 결과:

- readiness endpoint OK
- active targets: `13`
- targets up: `13`
- targets not up: `0`

표준 PrometheusRule sample을 실제 cluster에 적용했습니다.

적용한 파일:

- `rules/alerting/basic-alerts.yaml`
- `rules/prometheus/basic-recording-rules.yaml`
- `rules/slo/example-availability.yaml`

처음 발견한 문제:

- PrometheusRule CR은 생성됐지만 Prometheus API에 rule group이 로드되지 않았습니다.
- 원인은 installed `kube-prometheus-stack` Prometheus가 `release=kube-prometheus-stack` selector를 사용했는데, repository sample rule에 해당 label이 없었기 때문입니다.

수정:

- 표준 PrometheusRule sample에 kube-prometheus-stack selector용 label을 추가했습니다.

수정 후 Prometheus API에서 로드 확인된 rule group:

- `basic.kubernetes.alerting`
- `basic.kubernetes.recording`
- `slo.example-service.availability`

## Policy Sample 확인 결과

다음 policy sample들을 server-side dry-run으로 확인했습니다.

```powershell
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply --dry-run=server -f policies/network-policy-sample.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply --dry-run=server -f policies/pod-security-namespace.yaml
kubectl --kubeconfig .tmp/kubeconfig/local-k3s.yaml apply --dry-run=server -f policies/rbac-readonly-sample.yaml
```

결과:

- API server validation 통과

주의 사항:

- `pod-security-namespace.yaml`의 baseline enforce 적용 시 node-exporter pod에 대해 warning이 발생했습니다.
- 원인:
  - host namespace
  - hostPath volume
  - hostPort
- node-exporter DaemonSet 특성상 예상 가능한 warning으로 기록했습니다.

## GitHub 반영

작업 branch:

- `codex/traces-slo-profiles`

PR:

- `https://github.com/hkjs96/k8s-observability-standard/pull/2`

주요 commit:

- `4ba5737 Add traces and SLO profile validation`
- `37605f9 Add local k3d smoke path`
- `46392f1 Record local Basic smoke evidence`
- `f247b8c Document Grafana smoke credentials`

## 현재 상태

마지막 확인 기준:

- local k3d cluster `obs-standard-basic`는 살아 있습니다.
- monitoring pod는 6개 Running 상태입니다.
- Grafana는 port-forward를 통해 `http://127.0.0.1:3000`에서 접근 가능합니다.

## 남은 의사결정

1. Logs profile을 small deployment용 `SingleBinary`로 유지할지, production에서는 object storage 필수로 둘지 결정해야 합니다.
2. Traces profile에서 deprecated Tempo single binary chart를 계속 둘지, 다른 deployment model로 바꿀지 결정해야 합니다.
3. SLO generation 도구를 Sloth, Pyrra, provider-specific tooling 중 무엇으로 표준화할지 결정해야 합니다.
4. 실제 implementation repository adoption evidence를 별도로 수집해야 합니다.
