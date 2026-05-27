package smoke

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/example/k8s-observability/internal/validate"
)

type runner interface {
	Run(name string, args ...string) ([]byte, error)
	RunWithInputEnv(name string, args []string, input []byte, env []string) ([]byte, error)
}

type execRunner struct{}

func (execRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func (execRunner) RunWithInputEnv(name string, args []string, input []byte, env []string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), env...)
	if input != nil {
		cmd.Stdin = bytes.NewReader(input)
	}
	return cmd.CombinedOutput()
}

// Run executes smoke-test helper commands.
func Run(args []string) error {
	return run(args, execRunner{})
}

func run(args []string, r runner) error {
	if len(args) == 0 {
		return Usage()
	}

	switch args[0] {
	case "ec2-k3s":
		if len(args) < 2 {
			return Usage()
		}
		switch args[1] {
		case "launch":
			return launch(args[2:], r)
		case "fetch-kubeconfig":
			return fetchKubeconfig(args[2:], r)
		case "terminate":
			return terminate(args[2:], r)
		case "help", "-h", "--help":
			return Usage()
		default:
			return fmt.Errorf("unknown smoke ec2-k3s command %q", args[1])
		}
	case "k3s-basic":
		if len(args) < 2 {
			return Usage()
		}
		switch args[1] {
		case "install":
			return installK3sBasic(args[2:], r)
		case "help", "-h", "--help":
			return Usage()
		default:
			return fmt.Errorf("unknown smoke k3s-basic command %q", args[1])
		}
	case "help", "-h", "--help":
		return Usage()
	default:
		return fmt.Errorf("unknown smoke command %q", args[0])
	}
}

func launch(args []string, r runner) error {
	fs := flag.NewFlagSet("launch", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	region := fs.String("region", "", "AWS region")
	amiID := fs.String("ami-id", "", "Amazon Linux 2023 AMI ID")
	instanceType := fs.String("instance-type", "", "EC2 instance type")
	keyName := fs.String("key-name", "", "SSH key pair name")
	subnetID := fs.String("subnet-id", "", "subnet ID")
	securityGroupID := fs.String("security-group-id", "", "security group ID")
	name := fs.String("name", "example-k3s-basic-smoke", "instance Name tag")
	profile := fs.String("profile", "", "AWS CLI profile")
	userData := fs.String("user-data", filepath.Join("examples", "ec2-k3s-basic", "cloud-init.yaml"), "cloud-init file")
	volumeSize := fs.Int("volume-size-gib", 40, "root volume size in GiB")
	noWait := fs.Bool("no-wait", false, "skip waiting for instance status checks")
	yes := fs.Bool("yes", false, "actually launch the EC2 instance")
	dryRun := fs.Bool("dry-run", false, "print AWS CLI calls without executing")
	if err := fs.Parse(args); err != nil {
		return err
	}

	required := map[string]string{
		"--region":            *region,
		"--ami-id":            *amiID,
		"--instance-type":     *instanceType,
		"--key-name":          *keyName,
		"--subnet-id":         *subnetID,
		"--security-group-id": *securityGroupID,
	}
	if err := requireFlags(required); err != nil {
		return err
	}
	if _, err := os.Stat(*userData); err != nil {
		return fmt.Errorf("user data file not found: %s", *userData)
	}
	userDataPath, err := filepath.Abs(*userData)
	if err != nil {
		return err
	}
	if !*dryRun && !*yes {
		return errors.New("launch creates an EC2 instance; rerun with --yes or inspect with --dry-run")
	}

	blockDevices := []map[string]any{{
		"DeviceName": "/dev/xvda",
		"Ebs": map[string]any{
			"VolumeSize":          *volumeSize,
			"VolumeType":          "gp3",
			"DeleteOnTermination": true,
		},
	}}
	blockJSON, err := json.Marshal(blockDevices)
	if err != nil {
		return err
	}

	runArgs := []string{
		"ec2", "run-instances",
		"--region", *region,
		"--image-id", *amiID,
		"--instance-type", *instanceType,
		"--key-name", *keyName,
		"--subnet-id", *subnetID,
		"--security-group-ids", *securityGroupID,
		"--user-data", "file://" + filepath.ToSlash(userDataPath),
		"--metadata-options", "HttpTokens=required,HttpEndpoint=enabled",
		"--block-device-mappings", string(blockJSON),
		"--tag-specifications", fmt.Sprintf("ResourceType=instance,Tags=[{Key=Name,Value=%s},{Key=Purpose,Value=k3s-basic-smoke}]", *name),
		"--output", "json",
	}
	runArgs = appendProfile(runArgs, *profile)
	if *dryRun {
		printCommand("aws", runArgs)
		if !*noWait {
			printCommand("aws", appendProfile([]string{"ec2", "wait", "instance-status-ok", "--region", *region, "--instance-ids", "REPLACE_WITH_INSTANCE_ID"}, *profile))
		}
		return nil
	}

	out, err := r.Run("aws", runArgs...)
	if err != nil {
		return fmt.Errorf("aws ec2 run-instances failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	fmt.Print(string(out))

	instanceID, err := parseInstanceID(out)
	if err != nil {
		return err
	}
	if !*noWait {
		waitArgs := appendProfile([]string{"ec2", "wait", "instance-status-ok", "--region", *region, "--instance-ids", instanceID}, *profile)
		if out, err := r.Run("aws", waitArgs...); err != nil {
			return fmt.Errorf("aws ec2 wait instance-status-ok failed: %w\n%s", err, strings.TrimSpace(string(out)))
		}
	}
	describeArgs := appendProfile([]string{"ec2", "describe-instances", "--region", *region, "--instance-ids", instanceID, "--output", "json"}, *profile)
	out, err = r.Run("aws", describeArgs...)
	if err != nil {
		return fmt.Errorf("aws ec2 describe-instances failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	fmt.Print(string(out))
	return nil
}

func fetchKubeconfig(args []string, r runner) error {
	fs := flag.NewFlagSet("fetch-kubeconfig", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	host := fs.String("host", "", "EC2 host address")
	keyPath := fs.String("key-path", "", "SSH private key path")
	outputPath := fs.String("output", "", "kubeconfig output path")
	user := fs.String("user", "ec2-user", "SSH user")
	serverHost := fs.String("server-host", "", "kubeconfig API server host")
	port := fs.String("port", "6443", "kubeconfig API server port")
	force := fs.Bool("force", false, "overwrite existing output")
	allowRepoOutput := fs.Bool("allow-repository-output", false, "allow writing kubeconfig inside this repository")
	if err := fs.Parse(args); err != nil {
		return err
	}

	required := map[string]string{
		"--host":     *host,
		"--key-path": *keyPath,
		"--output":   *outputPath,
	}
	if err := requireFlags(required); err != nil {
		return err
	}
	if _, err := os.Stat(*keyPath); err != nil {
		return fmt.Errorf("ssh key not found: %s", *keyPath)
	}
	if _, err := os.Stat(*outputPath); err == nil && !*force {
		return fmt.Errorf("output already exists; rerun with --force: %s", *outputPath)
	}
	if err := rejectRepositoryOutput(*outputPath, *allowRepoOutput); err != nil {
		return err
	}
	if *serverHost == "" {
		*serverHost = *host
	}

	out, err := r.Run("ssh", "-i", *keyPath, "-o", "StrictHostKeyChecking=accept-new", *user+"@"+*host, "sudo cat /etc/rancher/k3s/k3s.yaml")
	if err != nil {
		return fmt.Errorf("ssh kubeconfig fetch failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}

	rendered := strings.ReplaceAll(string(out), "https://127.0.0.1:6443", "https://"+*serverHost+":"+*port)
	rendered = strings.ReplaceAll(rendered, "https://localhost:6443", "https://"+*serverHost+":"+*port)
	outputDir := filepath.Dir(*outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(*outputPath, []byte(rendered), 0o600); err != nil {
		return err
	}
	fmt.Println("wrote kubeconfig:", *outputPath)
	fmt.Println("do not commit this file if it contains a real endpoint")
	return nil
}

func terminate(args []string, r runner) error {
	fs := flag.NewFlagSet("terminate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	region := fs.String("region", "", "AWS region")
	instanceID := fs.String("instance-id", "", "EC2 instance ID")
	profile := fs.String("profile", "", "AWS CLI profile")
	noWait := fs.Bool("no-wait", false, "skip waiting for termination")
	yes := fs.Bool("yes", false, "actually terminate the EC2 instance")
	dryRun := fs.Bool("dry-run", false, "print AWS CLI calls without executing")
	if err := fs.Parse(args); err != nil {
		return err
	}
	required := map[string]string{
		"--region":      *region,
		"--instance-id": *instanceID,
	}
	if err := requireFlags(required); err != nil {
		return err
	}
	if !*dryRun && !*yes {
		return errors.New("terminate destroys an EC2 instance; rerun with --yes or inspect with --dry-run")
	}

	terminateArgs := appendProfile([]string{"ec2", "terminate-instances", "--region", *region, "--instance-ids", *instanceID, "--output", "json"}, *profile)
	if *dryRun {
		printCommand("aws", terminateArgs)
		if !*noWait {
			printCommand("aws", appendProfile([]string{"ec2", "wait", "instance-terminated", "--region", *region, "--instance-ids", *instanceID}, *profile))
		}
		return nil
	}
	out, err := r.Run("aws", terminateArgs...)
	if err != nil {
		return fmt.Errorf("aws ec2 terminate-instances failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	fmt.Print(string(out))
	if !*noWait {
		waitArgs := appendProfile([]string{"ec2", "wait", "instance-terminated", "--region", *region, "--instance-ids", *instanceID}, *profile)
		if out, err := r.Run("aws", waitArgs...); err != nil {
			return fmt.Errorf("aws ec2 wait instance-terminated failed: %w\n%s", err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func installK3sBasic(args []string, r runner) error {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	kubeconfig := fs.String("kubeconfig", "", "kubeconfig path")
	namespace := fs.String("namespace", "monitoring", "target namespace")
	releaseName := fs.String("release-name", "kube-prometheus-stack", "Helm release name")
	skipValidation := fs.Bool("skip-validation", false, "skip repository strict validation")
	dryRun := fs.Bool("dry-run", false, "print commands without executing")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *kubeconfig == "" {
		return errors.New("missing required flag: --kubeconfig")
	}
	if _, err := os.Stat(*kubeconfig); err != nil {
		return fmt.Errorf("kubeconfig not found: %s", *kubeconfig)
	}

	env := []string{"KUBECONFIG=" + *kubeconfig}
	commands := [][]string{
		{"kubectl", "get", "nodes"},
		{"helm", "repo", "add", "prometheus-community", "https://prometheus-community.github.io/helm-charts"},
		{"helm", "repo", "update"},
		{"helm", "upgrade", "--install", *releaseName, "prometheus-community/kube-prometheus-stack",
			"--version", "85.0.2",
			"--namespace", *namespace,
			"-f", "values/common/kube-prometheus-stack.yaml",
			"-f", "values/profiles/basic.yaml",
			"-f", "values/env/dev.yaml",
			"-f", "values/sizing/small.yaml",
			"--wait",
			"--timeout", "15m"},
		{"kubectl", "-n", *namespace, "get", "pods"},
		{"kubectl", "-n", *namespace, "get", "svc"},
		{"kubectl", "-n", *namespace, "get", "prometheusrules"},
		{"kubectl", "-n", *namespace, "get", "servicemonitors"},
		{"helm", "-n", *namespace, "status", *releaseName},
	}
	if *dryRun {
		if !*skipValidation {
			fmt.Println("go run ./cmd/obsctl validate --strict-tools")
		}
		printCommand("kubectl", []string{"create", "namespace", *namespace, "--dry-run=client", "-o", "yaml"})
		printCommand("kubectl", []string{"apply", "-f", "-"})
		printCommand("kubectl", []string{"-n", *namespace, "create", "secret", "generic", "grafana-admin", "--from-literal=admin-user=admin", "--from-literal=admin-password=REPLACE_FOR_SMOKE_ONLY", "--dry-run=client", "-o", "yaml"})
		printCommand("kubectl", []string{"apply", "-f", "-"})
		for _, command := range commands {
			printCommand(command[0], command[1:])
		}
		return nil
	}

	if !*skipValidation {
		if err := validate.Run("all", validate.Options{StrictTools: true}); err != nil {
			return err
		}
	}
	if err := runCommand(r, env, nil, "kubectl", "get", "nodes"); err != nil {
		return err
	}
	namespaceYAML, err := r.RunWithInputEnv("kubectl", []string{"create", "namespace", *namespace, "--dry-run=client", "-o", "yaml"}, nil, env)
	if err != nil {
		return fmt.Errorf("kubectl create namespace dry-run failed: %w\n%s", err, strings.TrimSpace(string(namespaceYAML)))
	}
	if err := runCommand(r, env, namespaceYAML, "kubectl", "apply", "-f", "-"); err != nil {
		return err
	}
	secretYAML, err := r.RunWithInputEnv("kubectl", []string{"-n", *namespace, "create", "secret", "generic", "grafana-admin", "--from-literal=admin-user=admin", "--from-literal=admin-password=REPLACE_FOR_SMOKE_ONLY", "--dry-run=client", "-o", "yaml"}, nil, env)
	if err != nil {
		return fmt.Errorf("kubectl create grafana-admin dry-run failed: %w\n%s", err, strings.TrimSpace(string(secretYAML)))
	}
	if err := runCommand(r, env, secretYAML, "kubectl", "apply", "-f", "-"); err != nil {
		return err
	}
	for _, command := range commands[1:] {
		if err := runCommand(r, env, nil, command[0], command[1:]...); err != nil {
			return err
		}
	}
	fmt.Println("k3s Basic smoke test ok")
	return nil
}

// Usage prints smoke helper usage.
func Usage() error {
	fmt.Println(`obsctl smoke helpers.

Usage:
  obsctl smoke ec2-k3s launch --region REGION --ami-id AMI --instance-type TYPE --key-name KEY --subnet-id SUBNET --security-group-id SG [--yes|--dry-run]
  obsctl smoke ec2-k3s fetch-kubeconfig --host HOST --key-path KEY --output PATH [--server-host HOST] [--force]
  obsctl smoke ec2-k3s terminate --region REGION --instance-id ID [--yes|--dry-run]
  obsctl smoke k3s-basic install --kubeconfig PATH [--namespace monitoring] [--release-name kube-prometheus-stack] [--dry-run]

Safety:
  launch and terminate require --yes for real AWS changes.
  use --dry-run to print the AWS CLI calls without executing them.
  fetch-kubeconfig refuses repository-local output unless --allow-repository-output is set.`)
	return nil
}

func appendProfile(args []string, profile string) []string {
	if profile == "" {
		return args
	}
	return append(args, "--profile", profile)
}

func requireFlags(values map[string]string) error {
	var missing []string
	for name, value := range values {
		if value == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required flag(s): %s", strings.Join(missing, ", "))
	}
	return nil
}

func parseInstanceID(raw []byte) (string, error) {
	var result struct {
		Instances []struct {
			InstanceID string `json:"InstanceId"`
		} `json:"Instances"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", err
	}
	if len(result.Instances) == 0 || result.Instances[0].InstanceID == "" {
		return "", errors.New("run-instances output did not include an instance ID")
	}
	return result.Instances[0].InstanceID, nil
}

func rejectRepositoryOutput(outputPath string, allow bool) error {
	if allow {
		return nil
	}
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	resolvedOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(repoRoot, resolvedOutput)
	if err != nil {
		return nil
	}
	if rel == "." || (!strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != "..") {
		return errors.New("output is inside this repository; use a path outside the repository or pass --allow-repository-output intentionally")
	}
	return nil
}

func printCommand(name string, args []string) {
	fmt.Println(shellJoin(append([]string{name}, args...)))
}

func runCommand(r runner, env []string, input []byte, name string, args ...string) error {
	out, err := r.RunWithInputEnv(name, args, input, env)
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	if err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}
	return nil
}

func shellJoin(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.ContainsAny(arg, " \t\n\"'[]{}") {
			quoted = append(quoted, strconvQuote(arg))
		} else {
			quoted = append(quoted, arg)
		}
	}
	return strings.Join(quoted, " ")
}

func strconvQuote(value string) string {
	escaped := strings.ReplaceAll(value, `"`, `\"`)
	return `"` + escaped + `"`
}
