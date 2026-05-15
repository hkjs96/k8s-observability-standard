---
id: pre-commit-validation
type: check
applies_to:
  - "**/*"
commands:
  - "go run ./cmd/obsctl validate"
---

# Pre-Commit Validation

Run:

```powershell
go run ./cmd/obsctl validate
```

Use this as the validation contract for local hooks or provider-specific CI.
