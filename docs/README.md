# Documentation

This directory is the human-facing standard documentation for Kubernetes
Observability v2.

Use these documents for:

- rationale and background
- standard decisions
- deployment and operations guidance
- security governance explanations
- validation contract and CI provider guidance
- Phase 0-2 completion status and known limits

Do not use `docs/` as the primary source for mandatory agent execution rules.
Agent rules, workflows, checks, and path routing live in `.agent/`.

## Index

- `00-overview.md`: repository purpose, boundary, and Phase 0-2 outcome
- `01-standard-decisions.md`: selected tools and deferred decisions
- `02-profiles.md`: profile boundaries
- `03-customization-guide.md`: values layering and implementation overrides
- `04-security-governance.md`: security guardrails and sensitive value handling
- `05-deployment-guide.md`: Basic profile deployment flow
- `06-day2-operations.md`: routine checks and operating assumptions
- `07-validation-contract.md`: portable validation commands and CI contract
- `08-phase-0-2-completion.md`: completed outputs, validation evidence, and known limits
