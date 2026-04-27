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
	rimpl "backend/src/internal/repository/impl"
	simpl "backend/src/internal/service/impl"
	"backend/src/internal/validator"

	sabst "backend/src/internal/service/abstract"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v3"
)

func Run() {
	config := config.Load()
	conn := postgres.NewPostgresConnection(config.GetDBDSN())
	serviceProvider := provider.NewServiceProvider()

	rosstatRepo := rimpl.NewRosstatRepository()
	rosstatAgeRepo := rimpl.NewRosstatAgeRepository()
	serviceProvider.Register((*sabst.IRosstatService)(nil), simpl.NewRosstatService(conn, rosstatRepo, rosstatAgeRepo))

	geoRepo := rimpl.NewGeoRepository()
	serviceProvider.Register((*sabst.IGeoService)(nil), simpl.NewGeoService(conn, geoRepo))

	app := fiber.New(fiber.Config{
		EnableSplittingOnParsers: true,
		StructValidator:          validator.NewFiberStructValidator(),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://midray.ru",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Accept, Content-Type, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/ping", health.PingHandler)

	app.Get("/api/v1/rosstat", middleware.Adapt(public.GetRosstatHandler, serviceProvider))
	app.Get("/api/v1/geo", middleware.Adapt(public.GetGeoHandler, serviceProvider))

	app.Get("/openapi.yaml", api.OpenapiYamlHandler)
	app.Get("/api/*", api.ApiHandler())

	log.Fatal(app.Listen(":80"))
}
