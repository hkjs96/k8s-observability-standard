---
id: sensitive-values-check
type: check
applies_to:
  - "**/*"
commands:
  - "go run ./cmd/obsctl validate sensitive"
---

# Sensitive Values Check

Run:

```powershell
go run ./cmd/obsctl validate sensitive
```

Completion criteria:

- No real credential-like values are committed.
- No forbidden Argo CD wildcard patterns are committed.
- `customers/` appears only in documentation that explicitly rejects that repository model.
