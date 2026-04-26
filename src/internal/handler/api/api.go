package api

import (
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

func OpenapiYamlHandler(c fiber.Ctx) error {
	return c.SendFile("src/api/public/openapi.yaml")
}

func ApiHandler() fiber.Handler {
	return swaggo.New(swaggo.Config{
		URL:         "/openapi.yaml",
		Title:       "API",
		DeepLinking: true,
	})
}
