# AGENTS.md

This repository uses shared agent instructions under `.agent/`.

Before editing:

1. Read `.agent/index.md`.
2. Read `.agent/context/repository-map.md`.
3. Read `.agent/rules/repository-boundary.md`.
4. Read `.agent/rules/sensitive-values.md`.
5. Select the relevant workflow from `.agent/workflows/` based on the task and changed paths.
6. Run the relevant checks from `.agent/checks/` before reporting completion.

Do not treat `docs/` as the source of mandatory execution rules. Use `docs/`
as human-facing reference material unless `.agent/index.md` explicitly points
to it.
