# CI Templates

CI provider choice is implementation-specific. This standard repository defines
the validation contract through `go run ./cmd/obsctl validate`, not through a
specific CI platform.

Use these files as examples when adapting the standard to a customer or internal
implementation repository:

- `github-actions-validate.yml`
- `gitlab-ci-validate.yml`
- `jenkinsfile-validate`

The CI job should validate, not deploy. Argo CD remains responsible for syncing
merged Git state to clusters.
