package services

import (
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
	artifacts := []models.Artifact{{Path: "generated/summary.md", Content: result}}
	_ = s.history.Save(userID, req.PipelineProvider, prompt)
	return artifacts, nil
}

func buildPrompt(req models.GenerationRequest) string {
	return fmt.Sprintf("Generate production CI/CD assets as YAML, Dockerfile, and K8s manifests for %s %s on %s for %s repository. Include security, scalability and monitoring (%s). Output JSON list of files.", req.Language, req.Framework, req.CloudProvider, req.RepositoryType, req.Monitoring)
}

func ParseGenerated(raw string) []models.Artifact {
	return []models.Artifact{{Path: "generated/output.txt", Content: strings.TrimSpace(raw)}}
}
