package user

import (
	"context"
  "fmt"
  "os"

	"github.com/gofiber/fiber/v2"

	"goweb/ui/forms"
  "goweb/db/user"
)

// login hanlder for fiber endpoint /login
// it expects a POST request
func Login(c *fiber.Ctx) error {
  creds := new(Credentials)
  if err := c.BodyParser(creds); err != nil {
    return err
  }

  ok, errs := ValidateCreds(creds)
  if ok == false {
    forms.Login(errs).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusBadRequest)
  }

  res, err := user.Get(creds.Username)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Query Failed: %v", err)
    return err
  }

  if res.Username == "" {
    errs["username"] = "username not found."
    forms.Login(errs).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusNotFound)
  }

  if creds.Password != res.Password {
    errs["password"] = "wrong password."
    forms.Login(errs).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusUnauthorized)
  }

  return c.SendString("You have successfully logged in 😄")
}

