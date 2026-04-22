package app

import (
	"log"

	"backend/src/internal/handler/api"
	"backend/src/internal/handler/health"
	"backend/src/internal/handler/public"
	"backend/src/internal/middleware"
	"backend/src/internal/provider"
	"backend/src/internal/service/abstract"
	"backend/src/internal/service/impl"
	"backend/src/internal/validator"

	"github.com/gofiber/fiber/v3"
)

func Run() {
	serviceProvider := provider.NewServiceProvider()
	serviceProvider.Register((*abstract.IRosstatService)(nil), &impl.RosstatService{})

	app := fiber.New(fiber.Config{
		EnableSplittingOnParsers: true,
		StructValidator:          validator.NewFiberStructValidator(),
	})

	app.Get("/ping", health.PingHandler)

	app.Get("/api/v1/rosstat", middleware.Adapt(public.RosstatHandler, serviceProvider))

	app.Get("/openapi.yaml", api.OpenapiYamlHandler)
	app.Get("/api/*", api.ApiHandler())

	log.Fatal(app.Listen(":8080"))
}
