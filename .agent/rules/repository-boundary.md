---
id: repository-boundary
type: rule
required: true
applies_to:
  - "**/*"
checks:
  - go run ./cmd/obsctl validate sensitive
---

# Repository Boundary Rules

- This is a standard repository, not an implementation repository.
- Do not create a `customers/` tree.
- Do not add real customer, account, cluster, domain, bucket, webhook, or IAM role values.
- Keep implementation-specific overrides in external implementation repositories.
- Use `examples/` only for fictional, non-deployable examples.
- Use `docs/` for human-facing explanations, not mandatory agent execution rules.
