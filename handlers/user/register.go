package user

import (
	"context"
  "fmt"
  "os"

	"github.com/gofiber/fiber/v2"

	"goweb/ui/forms"
  "goweb/db/user"
)

// register handle for fiber endpoint /register
// it expects a POST request
func Register(c *fiber.Ctx) error {
  creds := new(Credentials)
  if err := c.BodyParser(creds); err != nil {
    return err
  }

  ok, errs := ValidateCreds(creds)
  if ok == false {
    forms.Register(errs).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusBadRequest)
  }

  err := user.Add(creds.Username, creds.Password)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Query Failed: %v", err)
    return err
  }

  forms.Register(errs).Render(context.Background(), c.Response().BodyWriter())
  return c.SendStatus(fiber.StatusOK)
}
