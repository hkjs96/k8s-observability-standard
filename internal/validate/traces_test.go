package validate

import "testing"

func TestCheckTracePlaceholdersAcceptsProfileValues(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, root, "values/profiles/traces.yaml", "tempo:\n  retention: 24h\n")
	mustWrite(t, root, "values/profiles/traces-prometheus.yaml", "prometheus:\n  prometheusSpec:\n    enableRemoteWriteReceiver: true\n")
	mustWrite(t, root, "values/overrides/single-cluster-traces.yaml", "persistence:\n  enabled: true\n")
	mustWrite(t, root, "examples/opentelemetry/traces-instrumentation.yaml", "endpoint: http://tempo.observability-traces.svc.cluster.local:4317\n")
	mustWrite(t, root, "examples/phase3-smoke/example-application.yaml", "endpoint: tempo.observability-traces.svc.cluster.local:4317\n")
	mustWrite(t, root, "examples/phase3-smoke/trace-generator.yaml", "endpoint: tempo.observability-traces.svc.cluster.local:4317\n")
	withWorkingDir(t, root, func() {
		if err := checkTracePlaceholders(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCheckTracePlaceholdersRejectsImplementationValues(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, root, "values/profiles/traces.yaml", "endpoint: https://api.honeycomb.io\n")
	mustWrite(t, root, "values/profiles/traces-prometheus.yaml", "prometheus:\n  prometheusSpec:\n    enableRemoteWriteReceiver: true\n")
	mustWrite(t, root, "values/overrides/single-cluster-traces.yaml", "persistence:\n  enabled: true\n")
	mustWrite(t, root, "examples/opentelemetry/traces-instrumentation.yaml", "endpoint: http://tempo.observability-traces.svc.cluster.local:4317\n")
	mustWrite(t, root, "examples/phase3-smoke/example-application.yaml", "endpoint: tempo.observability-traces.svc.cluster.local:4317\n")
	mustWrite(t, root, "examples/phase3-smoke/trace-generator.yaml", "endpoint: tempo.observability-traces.svc.cluster.local:4317\n")
	withWorkingDir(t, root, func() {
		if err := checkTracePlaceholders(); err == nil {
			t.Fatal("checkTracePlaceholders() error = nil, want implementation value error")
		}
	})
}
