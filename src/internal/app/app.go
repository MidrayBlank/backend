package app

import (
	"log"

	"backend/src/internal/handler/health"
	"backend/src/internal/handler/public"
	"backend/src/internal/middleware"
	"backend/src/internal/provider"
	"backend/src/internal/service/abstract"
	"backend/src/internal/service/impl"

	"github.com/gofiber/fiber/v3"
)

func Run() {
	serviceProvider := provider.NewServiceProvider()
	serviceProvider.Register((*abstract.IRawService)(nil), &impl.RawService{})

	app := fiber.New(fiber.Config{
		EnableSplittingOnParsers: true,
	})

	app.Get("/", health.PingHandler)
	app.Get("/api/v1/rosstat/raw", middleware.Adapt(public.RosstatRawHandler, serviceProvider))

	log.Fatal(app.Listen(":8080"))
}
