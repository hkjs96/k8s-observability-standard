# Overview

This repository defines a repeatable Kubernetes observability standard. It
provides decisions, templates, profile values, policy samples, dashboard
candidates, and validation guidance that can be reused by implementation repos.

## Purpose

- Provide a standard baseline for Kubernetes monitoring with
  `kube-prometheus-stack`.
- Define profile boundaries for metrics, logs, traces, and SLO work.
- Keep common Helm values, Argo CD templates, rules, dashboards, and checklists
  in one versioned place.
- Make company-specific adoption predictable without storing company-specific
  runtime values here.

## Repository Boundary

This repository is not a central operations repository for many customer
clusters. It does not use a `customers/` directory. Real deployment values
belong in a separate implementation repository that references or copies this
standard after review.

Standard repository responsibilities:

- standard decisions and profiles
- common Helm values templates
- Argo CD Application and AppProject templates
- dashboard and rule templates
- security policy samples
- deployment and handover checklists

Implementation repository responsibilities:

- real cluster identifiers
- real endpoints and ingress hosts
- object storage settings
- IAM or workload identity bindings
- storage classes and retention values
- alert routes and notification secrets
- local dashboard and rule overrides
- operations runbooks

## Phase 0-2 Outcome

Phase 0-2 produces a usable Basic profile baseline:

- `kube-prometheus-stack` with embedded Grafana
- Prometheus Operator CRDs and core exporters
- Grafana dashboard sidecar convention
- PrometheusRule, ServiceMonitor, and PodMonitor ownership guidance
- Argo CD project and application templates using pinned chart versions
