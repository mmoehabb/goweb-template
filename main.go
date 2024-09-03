package main

import (
  "context"
	"github.com/gofiber/fiber/v2"

  "goweb/pages"
  "goweb/handlers/user"
)

func main() {
  // initialize a context to share data between different templ components
  ctx := context.WithValue(context.Background(), "version", "v0.0.1")

  app := fiber.New()
  app.Static("/public", "./public/")

  app.Get("/", func(c *fiber.Ctx) error {
    c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
    pages.Index().Render(ctx, c.Response().BodyWriter())
    return c.SendStatus(200)
  })

  app.Post("/login", user.Login)
  app.Post("/register", user.Register)

  app.Listen(":3000")
}
