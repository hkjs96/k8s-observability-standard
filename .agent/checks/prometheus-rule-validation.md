---
id: prometheus-rule-validation
type: check
applies_to:
  - "rules/**"
commands:
  - "go run ./cmd/obsctl validate prometheus"
---

# Prometheus Rule Validation

Run:

```powershell
go run ./cmd/obsctl validate prometheus
```

Completion criteria:

- Rule YAML parses.
- Each deployable rule sample has a matching `.promtool.yaml` mirror file.
- `promtool check rules` passes when promtool is available.
