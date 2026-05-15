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

import "gorm.io/gorm"

type DataModel struct {
    gorm.Model
    Field1 string `gorm:"uniqueIndex;not null"`
    Field2 string `gorm:"not null"`
}
```

### 2. Create Queries

`db/<feature>/queries.go`:

```go
package <feature>

import (
    "errors"

    "goweb/db"
)

func Add(field1, field2 string) error {
    var existing DataModel
    result := db.Connection.Where("field1 = ?", field1).First(&existing)
    if result.RowsAffected > 0 {
        return errors.New("already exists")
    }

    newModel := DataModel{Field1: field1, Field2: field2}
    result = db.Connection.Create(&newModel)
    return result.Error
}

func Get(field1 string) (DataModel, error) {
    var model DataModel
    result := db.Connection.Where("field1 = ?", field1).First(&model)
    if result.Error != nil {
        return DataModel{}, errors.New("not found")
    }
    return model, nil
}
```

### 3. Register in main.go

Add model to `RunMigrations()`:

```go
import "goweb/db/<feature>"

if err := db.RunMigrations(<feature>.DataModel{}); err != nil {
    log.Fatalf("Failed to run migrations: %v", err)
}
```

## Database Migrations

This template supports two migration strategies:

### 1. GORM AutoMigrate (Recommended for prototyping)

Models are defined with GORM tags and automatically synced to the database on startup.

**Running:** Migrations run automatically on `go run .` via `db.RunMigrations()` in `main.go`.

**Adding New Models:**
1. Define model in `db/<feature>/model.go` with GORM tags
2. Import the model in `main.go`
3. Add to `db.RunMigrations()` call

### 2. Goose SQL Migrations (Recommended for production)

For version-controlled SQL migrations with rollback support. SQL files go in `db/goose_migrations/`:

```sql
-- db/goose_migrations/001_create_<table>.sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS <table> (
    id SERIAL PRIMARY KEY,
    field1 VARCHAR(45) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS <table>;
-- +goose StatementEnd
```

**Running:** Also automatic on startup via `db.RunGooseMigrations()` in `main.go`.

**CLI Commands:**
```bash
goose postgres "postgres://user:pass@localhost:5432/db" up -dir ./db/goose_migrations
goose postgres "postgres://user:pass@localhost:5432/db" down -dir ./db/goose_migrations
```

Or use luci: `luci bash.migrate up`

### When to use which?

| Use Case | Recommended |
|----------|-------------|
| Fast prototyping | GORM AutoMigrate |
| Simple models | GORM AutoMigrate |
| Production apps | Goose |
| Complex SQL | Goose |
| Need rollbacks | Goose |

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

Use the global `db.Connection` for direct GORM operations:

```go
result := db.Connection.Where("field = ?", value).First(&model)
if result.Error != nil {
    return result.Error
}
```

Or use the wrapper for specific operations:

```go
conn, err := db.GetConnection()
if err != nil {
    return err
}
defer conn.Close()

result := conn.Create(&model)
return result.Error
```