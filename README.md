# doc-processing-project
A ready-to-run project with two Go microservices (Ingestion + Processor) demonstrating worker pools, Docker, Kubernetes manifests, and Terraform for provisioning Kubernetes resources.

## Overview
This project contains two microservices written in Go:
- **ingestion**: Accepts document upload requests at `/ingest`, enqueues jobs, and runs a worker pool that forwards jobs to the processor service.
- **processor**: Accepts jobs at `/process`, runs a worker pool to "process" documents (simulates heavy work), with retry/backoff.

## Running Locally with Docker Compose

1. Build and start both services:
   ```sh
   docker compose up --build
   ```

2. **Endpoints:**
   - Ingestion service: `POST http://localhost:8080/ingest`
   - Processor service: `POST http://localhost:8081/process`

3. **Environment Variables:**
   - The ingestion service uses the environment variable `PROCESSOR_URL` to forward jobs to the processor service. This is set in `docker-compose.yaml` as:
     ```yaml
     PROCESSOR_URL: "http://process:8081/process"
     ```

## Project Structure
- `ingestion/`: Source code and Dockerfile for the ingestion service
- `process/`: Source code and Dockerfile for the processor service
- `docker-compose.yaml`: Orchestrates both services

## Features
- Worker pool implementation in both services
- Graceful shutdown
- Containerized with multi-stage Dockerfiles
- Ready for deployment to Kubernetes

## Example Request
To enqueue a document for processing:
```sh
curl -X POST http://localhost:8080/ingest \
  -H "Content-Type: application/json" \
  -d '{"doc_id": "123", "content": "Hello World"}'
```

---

For more details, see the source code in each service directory.