package services

import "github.com/example/git-ops/backend/internal/repository"

type ArtifactService struct{ repo repository.ArtifactRepository }

func NewArtifactService(repo repository.ArtifactRepository) *ArtifactService { return &ArtifactService{repo: repo} }
