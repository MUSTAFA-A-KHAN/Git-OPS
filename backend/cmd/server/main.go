package main

import (
	"log"

	"github.com/example/git-ops/backend/internal/config"
	"github.com/example/git-ops/backend/internal/controllers"
	"github.com/example/git-ops/backend/internal/middleware"
	"github.com/example/git-ops/backend/internal/repository"
	"github.com/example/git-ops/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()
	db := repository.NewPostgres(cfg.DatabaseURL)
	artifactRepo := repository.NewArtifactRepository(db)
	genRepo := repository.NewGenerationRepository(db)
	authService := services.NewAuthService(cfg.JWTSecret)
	aiService := services.NewGenerationService(cfg, genRepo)
	artifactService := services.NewArtifactService(artifactRepo)

	app := fiber.New()
	app.Use(middleware.RequestLogger())

	authCtl := controllers.NewAuthController(authService)
	genCtl := controllers.NewGenerationController(aiService, artifactService)

	api := app.Group("/api/v1")
	authCtl.Register(api)
	gen := api.Group("/generate", middleware.JWT(cfg.JWTSecret))
	genCtl.Register(gen)

	log.Fatal(app.Listen(":" + cfg.Port))
}
