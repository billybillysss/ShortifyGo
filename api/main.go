package main

import (
	"github.com/billybillysss/ShortifyGo/database"
	"github.com/billybillysss/ShortifyGo/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/:short_id", handlers.RetrieveUrl)

	app.Post("/api/v1", handlers.ShortenUrl)

	app.Listen(":7001")

	defer database.RdsValid.Close()
	defer database.RdsUrl.Close()
}
