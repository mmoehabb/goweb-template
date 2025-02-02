package main

import (
	"context"

	"github.com/gofiber/fiber/v2"

	anc "goweb/ancillaries"
	"goweb/db"
	"goweb/handlers/user"
	"goweb/pages"
)

func main() {
	// initialize a context to share data between different templ components
	ctx := context.WithValue(context.Background(), "version", "v0.0.4")

	app := fiber.New()
	app.Static("/public", "./public/")

	// shall be used once and commented afterwards,
	// and maybe completed removed in production.
	app.Get("/seed", func(c *fiber.Ctx) error {
    defer anc.Recover(c)
		anc.Must(nil, db.Seed())
		return c.SendString("Database has been seeded.")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		pages.Index().Render(ctx, c.Response().BodyWriter())
		return c.SendStatus(200)
	})

	app.Post("/login", user.Login)
	app.Post("/register", user.Register)

	app.Listen(":3000")
}
