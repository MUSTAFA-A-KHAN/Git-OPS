package main

import (
	"encoding/json"
	"log"

	"github.com/example/git-ops/backend/internal/config"
	"github.com/example/git-ops/backend/internal/controllers"
	"github.com/example/git-ops/backend/internal/middleware"
	"github.com/example/git-ops/backend/internal/repository"
	"github.com/example/git-ops/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func getSwaggerJSON(scheme string, host string) []byte {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]string{
			"title":       "Git-OPS Backend API",
			"version":     "1.0.0",
			"description": "Swagger documentation for Git-OPS backend endpoints.",
		},
		"servers": []map[string]string{
			{"url": "/"},
		},
		"paths": map[string]interface{}{
			"/api/v1/auth/login": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Login and receive a demo token",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful login",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"token": map[string]string{"type": "string"},
										},
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/generate/": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Generate project artifacts",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"projectName":                  map[string]string{"type": "string"},
										"language":                     map[string]string{"type": "string"},
										"framework":                    map[string]string{"type": "string"},
										"cloudProvider":                map[string]string{"type": "string"},
										"deploymentTarget":             map[string]string{"type": "string"},
										"repositoryType":               map[string]string{"type": "string"},
										"testingFramework":             map[string]string{"type": "string"},
										"buildTool":                    map[string]string{"type": "string"},
										"containerizationPreference":   map[string]string{"type": "string"},
										"kubernetesUsage":              map[string]string{"type": "boolean"},
										"monitoringRequirements":       map[string]string{"type": "string"},
										"pipelineProvider":             map[string]string{"type": "string"},
									},
									"required": []string{"projectName"},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Generated artifacts",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"artifacts": map[string]interface{}{
												"type": "array",
												"items": map[string]interface{}{
													"type": "object",
													"properties": map[string]interface{}{
														"path":    map[string]string{"type": "string"},
														"content": map[string]string{"type": "string"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"security": []map[string]interface{}{
						{"bearerAuth": []string{}},
					},
				},
			},
			"/api/v1/generate/export": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Export the generated artifacts as a ZIP archive",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"artifacts": map[string]interface{}{
											"type": "array",
											"items": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"path":    map[string]string{"type": "string"},
													"content": map[string]string{"type": "string"},
												},
											},
										},
									},
									"required": []string{"artifacts"},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "ZIP archive of artifacts",
							"content": map[string]interface{}{
								"application/zip": map[string]interface{}{
									"schema": map[string]interface{}{
										"type":   "string",
										"format": "binary",
									},
								},
							},
						},
					},
					"security": []map[string]interface{}{
						{"bearerAuth": []string{}},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"bearerAuth": map[string]interface{}{
					"type":           "http",
					"scheme":         "bearer",
					"bearerFormat":   "JWT",
				},
			},
		},
	}

	data, _ := json.Marshal(spec)
	return data
}

var swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Git-OPS Backend Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      // Get the current origin (protocol + host)
      const origin = window.location.origin;
      
      SwaggerUIBundle({
        url: '/swagger.json',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: 'BaseLayout',
        validatorUrl: null,
        onComplete: function() {
          // After spec is loaded, override the servers to use current origin
          const spec = this.presets[0].lastSpec;
          if (spec && spec.servers) {
            spec.servers = [{ url: origin }];
          }
        },
        requestInterceptor: function(request) {
          // Make all requests relative to current origin
          if (!request.url.startsWith('http')) {
            request.url = origin + request.url;
          }
          return request;
        }
      });
    };
  </script>
</body>
</html>`

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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	authCtl := controllers.NewAuthController(authService)
	genCtl := controllers.NewGenerationController(aiService, artifactService)

	api := app.Group("/api/v1")
	authCtl.Register(api)
	gen := api.Group("/generate", middleware.JWT(cfg.JWTSecret))
	genCtl.Register(gen)

	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		scheme := "http"
		if c.Get("X-Forwarded-Proto") != "" {
			scheme = c.Get("X-Forwarded-Proto")
		}
		host := c.Get("Host")
		if host == "" {
			host = "localhost:8080"
		}
		c.Type("json")
		return c.Send(getSwaggerJSON(scheme, host))
	})
	app.Get("/swagger", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(swaggerHTML)
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}
