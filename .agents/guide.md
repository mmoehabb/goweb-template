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

### 1. Create Form Template

Create `ui/forms/<form>_templ.go` (write raw HTML/Templ):

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

Then run `templ generate` to generate the .go file.

### 2. Create Page Template

Create `pages/<page>_templ.go` and run `templ generate`.

Page files in `pages/` are auto-discovered by `ancillaries/endpoints.go` and mapped to routes:
- `pages/index_templ.go` → `/`
- `pages/about_templ.go` → `/about`
- `pages/user/profile_templ.go` → `/user/profile`

### 3. Reusable Components

Create in `ui/components/`:
- `Button_templ.go` → `Button(label, variant, attrs)`
- `TextInput_templ.go` → `TextInput(name, type, value, error, attrs)`

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

### 3. Add Seed (if new table)

In `db/db.go`, add to `Seed()`:

```go
"CREATE TABLE IF NOT EXISTS table (field1 VARCHAR(45) PRIMARY KEY, field2 VARCHAR(45) NOT NULL);",
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
| Form template | `ui/forms/login_templ.go` |
| Page template | `pages/index_templ.go` |
| DB model | `db/users/model.go` |
| DB queries | `db/users/queries.go` |

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