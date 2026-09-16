package main

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	anc "goweb/ancillaries"
	"goweb/constants"
	"goweb/db"
	"goweb/db/users"
	"goweb/handlers/user"
	"goweb/pages"
)

func main() {
	ctx := context.WithValue(context.Background(), "version", "v0.1.0")

	if _, err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.RunMigrations(users.DataModel{}); err != nil {
		log.Fatalf("Failed to run GORM migrations: %v", err)
	}

	app := fiber.New()
	app.Static("/public", "./public/")

	var endpoints = anc.GetEndpoint("./pages/")

	for _, endpoint := range endpoints {
		if page, ok := pages.Registry[endpoint]; ok {
			app.Get(endpoint, func(c *fiber.Ctx) error {
				c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
				page().Render(ctx, c.Response().BodyWriter())
				return c.SendStatus(200)
			})
		}
	}

	app.Post("/login", user.Login)
	app.Post("/register", user.Register)

	app.Listen(":" + strconv.Itoa(constants.AppConfig.Port))
}
