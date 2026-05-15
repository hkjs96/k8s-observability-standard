# Validation Contract

CI provider choice is implementation-specific. This standard does not require
GitHub Actions, GitLab CI, Jenkins, Azure DevOps, or any other CI product.

The portable validation contract is:

```bash
go run ./cmd/obsctl validate
```

For local Windows environments without Go, the fallback wrapper is:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/validate-all.ps1
```

## Responsibilities

Validation scripts check that the repository remains suitable for Argo CD
GitOps consumption:

- YAML files parse successfully.
- Basic profile values render with the pinned Helm chart when Helm is available.
- Argo CD templates avoid unsafe project and wildcard patterns.
- Prometheus rule samples have promtool mirror files.
- Sensitive or implementation-specific values are not committed.

## CI And Argo CD Boundary

CI validates Git state before merge. It must not deploy this standard directly
to customer clusters.

Argo CD owns sync and drift management after a reviewed Git state is merged into
an implementation repository.

Provider-specific examples are available under `templates/ci/`. Treat them as
starting points, not as the standard itself.
