# doc-processing-project
A ready-to-run small project with two Go microservices (Ingestion + Processor) that demonstrate worker pools, Docker, Kubernetes manifests, and Terraform to provision Kubernetes resources

Two microservices in Go showcasing worker pools and concurrency, deployable to Kubernetes and managed with Terraform.


Services:
- ingestion: accepts document upload requests and enqueues jobs; runs a worker pool that forwards jobs to processor.
- processor: accepts jobs and runs a worker pool to "process" documents (simulate heavy work), with retry/backoff.