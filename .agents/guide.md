# Developer Guide

This guide explains how to extend the GoWeb template with new features.

## Adding a New Feature

### 1. Create Handler

Create `handlers/<feature>/<handler>.go`:

```go
package <feature>

import (
    "context"

    "github.com/gofiber/fiber/v2"

    "goweb/db/<feature>"
    "goweb/ui/forms"
)

// <Feature> handles POST /<feature>
// expects request body matching your form
func <Feature>(c *fiber.Ctx) error {
    req := new(RequestStruct)
    if err := c.BodyParser(req); err != nil {
        return err
    }

    ok, errs := Validate(req)
    if !ok {
        forms.<FormName>(errs).Render(context.Background(), c.Response().BodyWriter())
        return c.SendStatus(fiber.StatusBadRequest)
    }

    err := <feature>.DoSomething(req)
    if err != nil {
        errs["field"] = err.Error()
        forms.<FormName>(errs).Render(context.Background(), c.Response().BodyWriter())
        return c.SendStatus(fiber.StatusBadRequest)
    }

    forms.<FormName>(errs).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusOK)
}
```

Handler naming: `FeatureName` (PascalCase, e.g., `Login`, `Register`, `CreatePost`)

### 2. Create Types (if needed)

Create `handlers/<feature>/types.go`:

```go
package <feature>

type RequestStruct struct {
    Field1 string `json:"field1" xml:"field1" form:"field1"`
    Field2 string `json:"field2" xml:"field2" form:"field2"`
}
```

### 3. Create Validators (optional)

Create `handlers/<feature>/validators.go`:

```go
package <feature>

func Validate(req *RequestStruct) (bool, map[string]string) {
    ok := true
    errs := make(map[string]string)

    if len(req.Field1) < 8 {
        errs["field1"] = "field1 should contain at least 8 characters."
        ok = false
    }

    return ok, errs
}
```

### 4. Add Route in main.go

```go
app.Post("/<feature>", <feature>.<Feature>)
```

## Creating UI Components

### Using templui (Recommended)

templui provides beautifully designed, accessible components via CLI workflow (like shadcn/ui).

**Setup**:
```bash
go install github.com/templui/templui/cmd/templui@latest
templui init
```

**Adding components**:
```bash
templui add button card dialog
```

**Usage**:
```templ
import "goweb/ui/components/button"

@button.Button(button.Props{Variant: button.VariantDefault, Type: button.TypeSubmit}) {
    Submit
}
```

**Available components**: Run `templui list` to see all available components.

**Updating**:
```bash
templui --installed add    # Update all installed components
templui upgrade            # Update CLI and utils
```

### Custom Components

Create custom components in `ui/components/` when templui doesn't have what you need.

```go
package forms

import "github.com/a-h/templ"

func Login(errs map[string]string) templ.Component {
    return templ.ComponentFunc(func(w io.Writer, _ templ.Context) error {
        _, err := io.WriteString(w, `<form hx-post="/login">`)
        // ... form HTML
        return err
    })
}
```

### Custom Form Template

Create `ui/forms/<form>.templ` (write raw HTML/Templ):

```go
package forms

import "github.com/a-h/templ"

templ Login(errs map[string]string) {
    <form hx-post="/login">
        // ... form HTML
    </form>
}
```

Then run `templ generate` to generate the .go file.

### 2. Create Page Template

Create `pages/<page>_templ.go` and run `templ generate`.

Page files in `pages/` are auto-discovered by `ancillaries/endpoints.go` and mapped to routes:
- `pages/index_templ.go` → `/`
- `pages/about_templ.go` → `/about`
- `pages/user/profile_templ.go` → `/user/profile`

### Reusable Components

Create in `ui/components/`:
- `TextInput.templ` → Custom input component (no templui equivalent yet)

## Database Layer

### 1. Create Model

`db/<feature>/model.go`:

```go
package <feature>

type DataModel struct {
    Field1 string
    Field2 string
}

func parseRow(row []any) DataModel {
    return DataModel{
        Field1: row[0].(string),
        Field2: row[1].(string),
    }
}
```

### 2. Create Queries

`db/<feature>/queries.go`:

```go
package <feature>

import (
    "errors"

    anc "goweb/ancillaries"
    "goweb/db"
)

func Add(field1, field2 string) error {
    conn := anc.Must(db.GetConnection()).(*db.Connection)

    res := anc.Must(conn.SeqQuery("SELECT * FROM table WHERE field1=$1", field1)).([]any)
    if len(res) != 0 {
        conn.Close()
        return errors.New("already exists")
    }

    anc.Must(conn.Query("INSERT INTO table VALUES ($1, $2)", field1, field2))
    return nil
}

func Get(field1 string) (DataModel, error) {
    conn := anc.Must(db.GetConnection()).(*db.Connection)

    res := anc.Must(conn.Query("SELECT * FROM table WHERE field1=$1", field1)).([]any)
    if len(res) == 0 {
        return DataModel{}, errors.New("not found")
    }

    return parseRow(res[0].([]any)), nil
}
```

## Database Migrations

Migrations are managed via [goose](https://github.com/pressly/goose). Create SQL migration files in `db/migrations/`:

### Creating a Migration

Create `db/migrations/XXX_create_<table>.sql`:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(45) REFERENCES users(username),
    bio TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profiles;
-- +goose StatementEnd
```

### Running Migrations

- **Auto:** Migrations run automatically on `go run .` (via `db.RunMigrations()`)
- **CLI:** `goose postgres "postgres://user:pass@localhost:5432/db?sslmode=disable" up -dir ./db/migrations`
- **Rollback:** `goose postgres "url" down -dir ./db/migrations`
- **Status:** `goose postgres "url" status -dir ./db/migrations`

### Creating New Migrations via CLI

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose create add_new_table sql
# Creates: db/migrations/YYYYMMDDHHMMSS_add_new_table.sql
```

## Conventions

| Pattern | Example |
|---------|---------|
| Handler file | `handlers/user/login.go` |
| Handler package | `package user` |
| Handler function | `func Login(c *fiber.Ctx) error` |
| Types file | `handlers/user/types.go` |
| Types struct | `type Credentials struct` |
| Validators | `handlers/user/validators.go` |
| Form template | `ui/forms/login.templ` |
| Page template | `pages/index.templ` |
| templui component | `ui/components/button/button.templ` |
| Custom UI component | `ui/components/TextInput.templ` |
| DB model | `db/users/model.go` |
| DB queries | `db/users/queries.go` |
| Migration file | `db/migrations/001_create_users.sql` |

### Validation Pattern

```go
func Validate*_(req *RequestType) (bool, map[string]string) {
    ok := true
    errs := make(map[string]string)

    if condition {
        errs["field"] = "error message"
        ok = false
    }

    return ok, errs
}
```

### Error Handling Pattern

```go
if !ok {
    forms.Form(errs).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusBadRequest)
}

if err != nil {
    errs["field"] = err.Error()
    forms.Form(errs).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusFound)
}
```

### Database Connection Pattern

```go
conn := anc.Must(db.GetConnection()).(*db.Connection)
defer conn.Close()

res := anc.Must(conn.Query(...)).([]any)
```

Where `anc.Must()` panics on error (defined in `ancillaries/`).