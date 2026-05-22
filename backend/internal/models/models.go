package models

import "time"

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type GenerationRequest struct {
	ProjectName          string `json:"projectName" validate:"required"`
	Language             string `json:"language"`
	Framework            string `json:"framework"`
	CloudProvider        string `json:"cloudProvider"`
	DeploymentTarget     string `json:"deploymentTarget"`
	RepositoryType       string `json:"repositoryType"`
	TestingFramework     string `json:"testingFramework"`
	BuildTool            string `json:"buildTool"`
	Containerization     string `json:"containerizationPreference"`
	KubernetesUsage      bool   `json:"kubernetesUsage"`
	Monitoring           string `json:"monitoringRequirements"`
	PipelineProvider     string `json:"pipelineProvider"`
}

type Artifact struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type Generation struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	UserID     string     `json:"userId"`
	Input      string     `json:"input"`
	Provider   string     `json:"provider"`
	Artifacts  []Artifact `json:"artifacts" gorm:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
}
