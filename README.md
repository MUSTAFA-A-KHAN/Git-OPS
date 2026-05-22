# Git-Ops AI Pipeline Generator

Enterprise-ready full-stack platform to generate CI/CD artifacts using **local Llama 3 via Ollama**.

## Stack
- Frontend: React + TypeScript + Vite
- Backend: Go + Fiber
- AI: Ollama REST API
- Infra: Docker Compose, Kubernetes

## Features
- Pipeline generation (GitHub Actions, GitLab CI, Jenkins, Azure DevOps)
- Dockerfile + K8s manifests generation
- Deployment and infrastructure recommendations
- JWT auth, generation history, ZIP export

## API
- `POST /api/v1/auth/login`
- `POST /api/v1/generate/`
- `POST /api/v1/generate/export`

## Setup
```bash
docker compose up --build
```
