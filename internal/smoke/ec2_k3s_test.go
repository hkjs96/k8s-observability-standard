package smoke

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	calls [][]string
	out   []byte
	err   error
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	call := append([]string{name}, args...)
	f.calls = append(f.calls, call)
	return f.out, f.err
}

func (f *fakeRunner) RunWithInputEnv(name string, args []string, input []byte, env []string) ([]byte, error) {
	call := append([]string{name}, args...)
	f.calls = append(f.calls, call)
	return f.out, f.err
}

func TestLaunchRequiresYesUnlessDryRun(t *testing.T) {
	userData := filepath.Join(t.TempDir(), "cloud-init.yaml")
	if err := os.WriteFile(userData, []byte("#cloud-config\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := run([]string{
		"ec2-k3s", "launch",
		"--region", "example-region-1",
		"--ami-id", "ami-example",
		"--instance-type", "t3.large",
		"--key-name", "example-key",
		"--subnet-id", "subnet-example",
		"--security-group-id", "sg-example",
		"--user-data", userData,
	}, &fakeRunner{})
	if err == nil || !strings.Contains(err.Error(), "--yes") {
		t.Fatalf("launch error = %v, want --yes safety error", err)
	}
}

func TestLaunchDryRunDoesNotCallRunner(t *testing.T) {
	userData := filepath.Join(t.TempDir(), "cloud-init.yaml")
	if err := os.WriteFile(userData, []byte("#cloud-config\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{}
	err := run([]string{
		"ec2-k3s", "launch",
		"--region", "example-region-1",
		"--ami-id", "ami-example",
		"--instance-type", "t3.large",
		"--key-name", "example-key",
		"--subnet-id", "subnet-example",
		"--security-group-id", "sg-example",
		"--user-data", userData,
		"--dry-run",
	}, runner)
	if err != nil {
		t.Fatal(err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("runner calls = %d, want 0", len(runner.calls))
	}
}

func TestTerminateRequiresYesUnlessDryRun(t *testing.T) {
	err := run([]string{
		"ec2-k3s", "terminate",
		"--region", "example-region-1",
		"--instance-id", "i-example",
	}, &fakeRunner{})
	if err == nil || !strings.Contains(err.Error(), "--yes") {
		t.Fatalf("terminate error = %v, want --yes safety error", err)
	}
}

func TestFetchKubeconfigRejectsRepositoryOutput(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(keyPath, []byte("example-key"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := run([]string{
		"ec2-k3s", "fetch-kubeconfig",
		"--host", "example.invalid",
		"--key-path", keyPath,
		"--output", filepath.Join("tmp-kubeconfig.yaml"),
	}, &fakeRunner{})
	if err == nil || !strings.Contains(err.Error(), "inside this repository") {
		t.Fatalf("fetch error = %v, want repository output rejection", err)
	}
}

func TestFetchKubeconfigWritesOutsideRepository(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	if err := os.WriteFile(keyPath, []byte("example-key"), 0o600); err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(dir, "kubeconfig.yaml")
	runner := &fakeRunner{
		out: []byte("clusters:\n- cluster:\n    server: https://127.0.0.1:6443\n"),
	}

	err := run([]string{
		"ec2-k3s", "fetch-kubeconfig",
		"--host", "example.invalid",
		"--key-path", keyPath,
		"--output", outputPath,
	}, runner)
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rendered), "https://example.invalid:6443") {
		t.Fatalf("rendered kubeconfig = %q, want rewritten server", string(rendered))
	}
}

func TestParseInstanceID(t *testing.T) {
	id, err := parseInstanceID([]byte(`{"Instances":[{"InstanceId":"i-example"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if id != "i-example" {
		t.Fatalf("id = %q, want i-example", id)
	}
}

func TestRunnerErrorIncludesOutput(t *testing.T) {
	userData := filepath.Join(t.TempDir(), "cloud-init.yaml")
	if err := os.WriteFile(userData, []byte("#cloud-config\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := run([]string{
		"ec2-k3s", "launch",
		"--region", "example-region-1",
		"--ami-id", "ami-example",
		"--instance-type", "t3.large",
		"--key-name", "example-key",
		"--subnet-id", "subnet-example",
		"--security-group-id", "sg-example",
		"--user-data", userData,
		"--yes",
	}, &fakeRunner{out: []byte("example failure"), err: errors.New("failed")})
	if err == nil || !strings.Contains(err.Error(), "example failure") {
		t.Fatalf("error = %v, want command output", err)
	}
}

func TestInstallK3sBasicDryRunDoesNotCallRunner(t *testing.T) {
	kubeconfig := filepath.Join(t.TempDir(), "kubeconfig.yaml")
	if err := os.WriteFile(kubeconfig, []byte("apiVersion: v1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{}
	err := run([]string{
		"k3s-basic", "install",
		"--kubeconfig", kubeconfig,
		"--dry-run",
	}, runner)
	if err != nil {
		t.Fatal(err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("runner calls = %d, want 0", len(runner.calls))
	}
}
