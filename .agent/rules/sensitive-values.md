---
id: sensitive-values
type: rule
required: true
applies_to:
  - "**/*"
checks:
  - go run ./cmd/obsctl validate sensitive
---

# Sensitive Values Rules

- Do not commit plaintext credentials.
- Do not commit real customer names, account IDs, bucket names, hostnames, webhook URLs, IAM role ARNs, or cluster names.
- Use placeholders such as `example-*`, `*.example.invalid`, or `REPLACE_IN_IMPLEMENTATION_REPO`.
- Grafana credentials must be referenced through an existing Secret.
- Alert receivers and notification routes must be provided by implementation repositories.
