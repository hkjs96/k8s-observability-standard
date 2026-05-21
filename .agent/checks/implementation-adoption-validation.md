---
id: implementation-adoption-validation
type: check
applies_to:
  - "docs/09-implementation-adoption.md"
  - "templates/implementation-adoption-checklist.template.md"
  - "examples/**"
commands:
  - "go run ./cmd/obsctl validate --strict-tools"
  - "go run ./cmd/obsctl validate sensitive"
---

# Implementation Adoption Validation

Run:

```powershell
go run ./cmd/obsctl validate --strict-tools
go run ./cmd/obsctl validate sensitive
```

Completion criteria:

- Adoption guidance keeps implementation-owned runtime values outside this repository.
- Examples remain fictional and non-deployable.
- Checklist fields use placeholders, not real cluster, account, domain, bucket, webhook, or identity values.
- Strict validation and sensitive scan remain green.
