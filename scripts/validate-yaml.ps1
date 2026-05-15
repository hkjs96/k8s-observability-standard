param(
  [string[]]$Paths = @(".")
)

$ErrorActionPreference = "Stop"

$python = Get-Command python -ErrorAction SilentlyContinue
if (-not $python) {
  Write-Error "python is required for YAML validation in this repository."
  exit 1
}

$joined = ($Paths | ForEach-Object { $_ }) -join [System.IO.Path]::PathSeparator

$script = @'
import pathlib
import sys

try:
    import yaml
except Exception as exc:
    print(f"PyYAML unavailable: {exc}", file=sys.stderr)
    sys.exit(2)

roots = sys.argv[1].split(__import__("os").pathsep)
errors = []
for root in roots:
    path = pathlib.Path(root)
    candidates = [path] if path.is_file() else path.rglob("*.yaml")
    for candidate in candidates:
        parts = set(candidate.parts)
        if ".cache" in parts:
            continue
        if candidate.suffix not in {".yaml", ".yml"}:
            continue
        try:
            text = candidate.read_text(encoding="utf-8")
            list(yaml.safe_load_all(text))
        except Exception as exc:
            errors.append(f"{candidate}: {exc}")

if errors:
    print("\n".join(errors), file=sys.stderr)
    sys.exit(1)

print("yaml ok")
'@

$script | python - $joined
