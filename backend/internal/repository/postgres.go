package repository

import "gorm.io/gorm"

type GenerationRepository interface{ Save(userID, provider, input string) error }
type ArtifactRepository interface{}

type generationRepo struct{ db *gorm.DB }
type artifactRepo struct{ db *gorm.DB }

func NewPostgres(_ string) *gorm.DB { return &gorm.DB{} }
func NewGenerationRepository(db *gorm.DB) GenerationRepository { return &generationRepo{db: db} }
func NewArtifactRepository(db *gorm.DB) ArtifactRepository { return &artifactRepo{db: db} }
func (r *generationRepo) Save(userID, provider, input string) error { return nil }
