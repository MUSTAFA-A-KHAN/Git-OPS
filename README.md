# Git-Ops AI Pipeline Generator

Enterprise-ready full-stack platform to generate CI/CD artifacts using an **OpenAI-compatible chat completions API**.

## Stack
- Frontend: React + TypeScript + Vite
- Backend: Go + Fiber
- AI: OpenAI-compatible Chat Completions API
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
set OPENAI_API_KEY=your-api-key
docker compose up --build
```
