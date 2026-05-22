package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/example/git-ops/backend/internal/ai"
	"github.com/example/git-ops/backend/internal/config"
	"github.com/example/git-ops/backend/internal/models"
	"github.com/example/git-ops/backend/internal/repository"
)

type GenerationService struct {
	cfg     config.Config
	client  ai.Client
	history repository.GenerationRepository
}

func NewGenerationService(cfg config.Config, history repository.GenerationRepository) *GenerationService {
	return &GenerationService{
		cfg:     cfg,
		client:  ai.Client{BaseURL: cfg.OpenAIURL, APIKey: cfg.OpenAIKey, Model: cfg.OpenAIModel},
		history: history,
	}
}

func (s *GenerationService) Generate(userID string, req models.GenerationRequest) ([]models.Artifact, error) {
	prompt := buildPrompt(req)
	result, err := s.client.Generate(prompt)
	if err != nil {
		return nil, err
	}
	artifacts := ParseGenerated(result)
	_ = s.history.Save(userID, req.PipelineProvider, prompt)
	return artifacts, nil
}

func buildPrompt(req models.GenerationRequest) string {
	return fmt.Sprintf(`Generate production CI/CD assets for this project:
- Project name: %s
- Language/framework: %s/%s
- Cloud provider: %s
- Deployment target: %s
- Repository type: %s
- Testing framework: %s
- Build tool: %s
- Containerization: %s
- Kubernetes manifests: %t
- Monitoring requirements: %s
- Pipeline provider: %s

Return only valid JSON. Do not wrap it in Markdown fences.
The JSON must be an array of file objects using this exact shape:
[{"path":"relative/file/path","content":"complete file content"}]

Include production-ready pipeline YAML, Dockerfile, Kubernetes manifests when requested, and a short README explaining how to use the generated files.`,
		req.ProjectName,
		req.Language,
		req.Framework,
		req.CloudProvider,
		req.DeploymentTarget,
		req.RepositoryType,
		req.TestingFramework,
		req.BuildTool,
		req.Containerization,
		req.KubernetesUsage,
		req.Monitoring,
		req.PipelineProvider,
	)
}

func ParseGenerated(raw string) []models.Artifact {
	cleaned := stripMarkdownFence(strings.TrimSpace(raw))
	if cleaned == "" {
		return []models.Artifact{{Path: "generated/output.txt", Content: ""}}
	}

	var artifacts []models.Artifact
	if err := json.Unmarshal([]byte(cleaned), &artifacts); err == nil && len(artifacts) > 0 {
		return sanitizeArtifacts(artifacts)
	}

	var wrapped struct {
		Artifacts []models.Artifact `json:"artifacts"`
		Files     []models.Artifact `json:"files"`
	}
	if err := json.Unmarshal([]byte(cleaned), &wrapped); err == nil {
		if len(wrapped.Artifacts) > 0 {
			return sanitizeArtifacts(wrapped.Artifacts)
		}
		if len(wrapped.Files) > 0 {
			return sanitizeArtifacts(wrapped.Files)
		}
	}

	return []models.Artifact{{Path: "generated/output.md", Content: cleaned}}
}

func stripMarkdownFence(value string) string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "```") {
		return value
	}

	value = strings.TrimPrefix(value, "```")
	if newline := strings.Index(value, "\n"); newline >= 0 {
		value = value[newline+1:]
	}
	return strings.TrimSpace(strings.TrimSuffix(value, "```"))
}

func sanitizeArtifacts(artifacts []models.Artifact) []models.Artifact {
	cleaned := make([]models.Artifact, 0, len(artifacts))
	for i, artifact := range artifacts {
		path := strings.TrimSpace(artifact.Path)
		if path == "" {
			path = fmt.Sprintf("generated/file-%d.txt", i+1)
		}
		cleaned = append(cleaned, models.Artifact{
			Path:    path,
			Content: strings.TrimSpace(artifact.Content),
		})
	}
	return cleaned
}
