package controllers

import (
	"archive/zip"
	"bytes"

	"github.com/example/git-ops/backend/internal/models"
	"github.com/example/git-ops/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type GenerationController struct {
	service   *services.GenerationService
	artifacts *services.ArtifactService
}

func NewGenerationController(s *services.GenerationService, a *services.ArtifactService) *GenerationController { return &GenerationController{service: s, artifacts: a} }
func (g *GenerationController) Register(r fiber.Router) { r.Post("/", g.Generate); r.Post("/export", g.ExportZip) }

func (g *GenerationController) Generate(c *fiber.Ctx) error {
	var req models.GenerationRequest
	if err := c.BodyParser(&req); err != nil { return fiber.ErrBadRequest }
	files, err := g.service.Generate("demo-user", req)
	if err != nil { return fiber.NewError(fiber.StatusBadGateway, err.Error()) }
	return c.JSON(fiber.Map{"artifacts": files})
}

func (g *GenerationController) ExportZip(c *fiber.Ctx) error {
	var payload struct{ Artifacts []models.Artifact `json:"artifacts"` }
	if err := c.BodyParser(&payload); err != nil { return fiber.ErrBadRequest }
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, a := range payload.Artifacts {
		f, _ := zw.Create(a.Path)
		_, _ = f.Write([]byte(a.Content))
	}
	_ = zw.Close()
	c.Set("Content-Type", "application/zip")
	return c.Send(buf.Bytes())
}
