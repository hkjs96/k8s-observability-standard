package validate

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/example/k8s-observability/internal/walk"
)

func YAML() error {
	files, err := walk.Files(".", func(path string) bool {
		return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
	})
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	python, err := exec.LookPath("python")
	if err != nil {
		fmt.Println("python unavailable; skipped YAML parser check")
		return nil
	}

	payload, err := json.Marshal(files)
	if err != nil {
		return err
	}

	code := `
import json
import pathlib
import sys
try:
    import yaml
except Exception as exc:
    print(f"PyYAML unavailable; skipped YAML parser check: {exc}")
    sys.exit(0)
errors = []
for raw in json.loads(sys.argv[1]):
    path = pathlib.Path(raw)
    try:
        list(yaml.safe_load_all(path.read_text(encoding="utf-8")))
    except Exception as exc:
        errors.append(f"{raw}: {exc}")
if errors:
    print("\n".join(errors), file=sys.stderr)
    sys.exit(1)
print("yaml ok")
`

	out, err := exec.Command(python, "-c", code, string(payload)).CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	if err != nil {
		return fmt.Errorf("YAML validation failed: %w", err)
	}
	return nil
}
