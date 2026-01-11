package main

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/tools/go/packages"

	anc "goweb/ancillaries"
	"goweb/constants"
	"goweb/db"
	"goweb/handlers/user"
	"goweb/pages"
)

func main() {
	// initialize a context to share data between different templ components
	ctx := context.WithValue(context.Background(), "version", "v0.0.6")

	app := fiber.New()
	app.Static("/public", "./public/")

	// shall be used once and commented afterwards,
	// and maybe completed removed in production.
	app.Get("/seed", func(c *fiber.Ctx) error {
		defer anc.Recover(c)
		anc.Must(nil, db.Seed())
		return c.SendString("Database has been seeded.")
	})

	var endpoints = anc.GetEndpoint("./pages/")

	for _, endpoint := range endpoints {
		pkgs, err := packages.Load(&packages.Config{}, "./pages/"+endpoint)
		if err != nil {
			log.Println(err)
			continue
		}
		packages.Visit(pkgs, nil, func(p *packages.Package) {
		})

		app.Get(endpoint, func(c *fiber.Ctx) error {
			c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
			pages.Index().Render(ctx, c.Response().BodyWriter())
			return c.SendStatus(200)
		})
	}

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		pages.Index().Render(ctx, c.Response().BodyWriter())
		return c.SendStatus(200)
	})

	app.Post("/login", user.Login)
	app.Post("/register", user.Register)

	app.Listen(":" + strconv.Itoa(constants.AppConfig.Port))
}
