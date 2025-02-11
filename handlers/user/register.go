package user

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"goweb/db/users"
	"goweb/ui/forms"
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

	err := users.Add(creds.Username, creds.Password)
	if err != nil {
		errs["username"] = err.Error()
		forms.Register(errs).Render(context.Background(), c.Response().BodyWriter())
		return c.SendStatus(fiber.StatusFound)
	}

	forms.Register(errs).Render(context.Background(), c.Response().BodyWriter())
	return c.SendStatus(fiber.StatusOK)
}
