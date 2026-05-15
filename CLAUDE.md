# CLAUDE.md

Use `.agent/` as the primary shared instruction source for this repository.

Required reading before edits:

1. `.agent/index.md`
2. `.agent/context/repository-map.md`
3. `.agent/rules/repository-boundary.md`
4. `.agent/rules/sensitive-values.md`

For task-specific work, follow the matching workflow in `.agent/workflows/`.

If a workflow requires validation, run the commands described in `.agent/checks/`.
If a rule conflicts with user instructions, stop and explain the conflict before editing.
