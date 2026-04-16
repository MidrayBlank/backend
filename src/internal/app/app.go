package app

import (
	"log"

	"backend/src/internal/handler/health"

	"github.com/gofiber/fiber/v3"
)

func Run() {
	app := fiber.New()

	app.Get("/", health.PingHandler)

	log.Fatal(app.Listen(":8080"))
}
