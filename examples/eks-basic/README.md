# EKS Basic Example

This example is for planning an EKS implementation of the Basic profile. It does
not include real AWS account values.

Implementation repo responsibilities:

- choose the final namespace name
- provide workload identity or IRSA annotations if needed
- select storage class for Prometheus persistence
- decide whether EKS control plane metrics and logs need CloudWatch integration
- provide Grafana access through an approved private path

Keep AWS account IDs, role ARNs, bucket names, and endpoint hosts out of this
standard repository.
