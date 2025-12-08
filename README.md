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

## Running on Kubernetes

1. **Apply resources:**
   ```sh
   kubectl apply -f k8s/namespace.yaml
   kubectl apply -f k8s/configmap.yaml
   kubectl apply -f k8s/ingestion.yaml
   kubectl apply -f k8s/process-deployment.yaml
   kubectl apply -f k8s/ingestion-service.yaml
   kubectl apply -f k8s/process-service.yaml
   kubectl apply -f k8s/ingestion-hpa.yaml
   kubectl apply -f k8s/process-hpa.yaml
   ```

2. **Endpoints:**
   - Ingestion service: Exposed via NodePort (default: `30080`). Access with `http://<NodeIP>:30080/ingest`
   - Processor service: Exposed via NodePort (default: `30081`). Access with `http://<NodeIP>:30081/process`

3. **Scaling:**
   - Horizontal Pod Autoscaler (HPA) automatically scales pods based on CPU usage.

4. **Health Checks:**
   - Both services implement `/health` endpoints for readiness and liveness probes.

5. **Config & Environment:**
   - Configuration is managed via ConfigMaps and environment variables in manifests.

## Project Structure
- `ingestion/`: Source code and Dockerfile for the ingestion service
- `process/`: Source code and Dockerfile for the processor service
- `k8s/`: Kubernetes manifests for deployments, services, HPA, configmaps, and namespace
- `docker-compose.yaml`: Orchestrates both services for local development

## Features
- Worker pool implementation in both services
- Graceful shutdown
- Containerized with multi-stage Dockerfiles
- Ready for deployment to Kubernetes
- Horizontal Pod Autoscaling
- Health checks for robust operation

## Example Request
To enqueue a document for processing:
```sh
curl -X POST http://localhost:8080/ingest \
  -H "Content-Type: application/json" \
  -d '{"doc_id": "123", "content": "Hello World"}'
```

For Kubernetes NodePort, replace `localhost:8080` with `<NodeIP>:30080`.

---

For more details, see the source code in each service directory and the manifests in the `k8s/` folder.