# Templates

Templates are reusable starting points for implementation repositories and
handover workflows. They must stay generic and must not include runtime-specific
values.

## Index

- `assessment-sheet.template.md`: initial observability assessment.
- `basic-smoke-evidence.template.md`: evidence capture for disposable Basic profile smoke tests.
- `deployment-checklist.template.md`: Basic profile deployment checklist.
- `implementation-adoption-checklist.template.md`: implementation repository adoption checklist.
- `operation-handover.template.md`: operations handover skeleton.
- `security-review.template.md`: security review checklist.
- `ci/`: CI validation examples.

## Rules

- Keep values fictional or blank.
- Do not add real secrets, endpoints, account IDs, bucket names, domains, webhook URLs, or identity bindings.
- Keep deployment-specific overrides in implementation repositories.
- Use `docs/09-implementation-adoption.md` for adoption guidance.
