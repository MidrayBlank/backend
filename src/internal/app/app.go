package app

import (
	"log"

	"backend/src/internal/config"
	"backend/src/internal/db/postgres"
	"backend/src/internal/handler/api"
	"backend/src/internal/handler/health"
	"backend/src/internal/handler/public"
	"backend/src/internal/middleware"
	"backend/src/internal/provider"
	sabst "backend/src/internal/service/abstract"
	simpl "backend/src/internal/service/impl"
	"backend/src/internal/validator"

	"github.com/gofiber/fiber/v3"
)

func Run() {
	config := config.Load()
	conn := postgres.NewPostgresConnection(config.GetDBDSN())

	serviceProvider := provider.NewServiceProvider()
	serviceProvider.Register((*sabst.IRosstatService)(nil), simpl.NewRosstatService(conn))

	app := fiber.New(fiber.Config{
		EnableSplittingOnParsers: true,
		StructValidator:          validator.NewFiberStructValidator(),
	})

	app.Get("/ping", health.PingHandler)

	app.Get("/api/v1/rosstat", middleware.Adapt(public.GetRosstatInfo, serviceProvider))
	//app.Get("/api/v1/geo", middleware.Adapt(public.GetHandlerInfo, serviceProvider))

	app.Get("/openapi.yaml", api.OpenapiYamlHandler)
	app.Get("/api/*", api.ApiHandler())

	log.Fatal(app.Listen(":80"))
}
