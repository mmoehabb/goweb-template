# Tests

This directory contains unit tests for the goweb-template application.

## Running Tests

```bash
# Run all tests
go test ./tests/...

# Run tests for specific package
go test ./tests/handlers/user/...
go test ./tests/db/users/...
go test ./tests/ancillaries/...
```

## Test Structure

```
tests/
├── handlers/
│   └── user/
│       └── validators_test.go    # Tests for user validation logic
├── db/
│   └── users/
│       └── queries_test.go        # Tests for database operations
├── ancillaries/
│   └── endpoints_test.go          # Tests for endpoint discovery
└── README.md                      # This file
```

## Writing Tests

### Handler Tests

Tests for handler logic (validators, types) go in `tests/handlers/<feature>/`.

Example:
```go
package user

import (
    "testing"

    "goweb/handlers/user"
)

func TestValidateCreds_Valid(t *testing.T) {
    creds := &user.Credentials{
        Username: "testuser",
        Password: "testpassword",
    }

    ok, errs := user.ValidateCreds(creds)

    if !ok {
        t.Errorf("expected valid credentials, got errors: %v", errs)
    }
}
```

### Database Tests

Database tests use SQLite in-memory for testing without requiring a PostgreSQL instance.

The `db/users/queries.go` provides test helper functions:
- `SetConnectionForTest(conn *gorm.DB)` - Set a test database connection
- `ResetConnectionForTest()` - Reset to original connection

Example:
```go
func TestAdd_CreatesUser(t *testing.T) {
    db := setupTestDB(t) // Creates SQLite in-memory DB

    originalConnection := users.SetConnectionForTest(db)
    defer users.ResetConnectionForTest()

    err := users.Add("testuser", "testpassword")
    // ...
}
```

### Ancillary Tests

Tests for utility functions (endpoints, errors) go in `tests/ancillaries/`.

## Test Dependencies

- Go's built-in `testing` package
- `gorm.io/driver/sqlite` - For database testing without PostgreSQL

## Notes

- Tests use Go's standard `testing` package (no external test frameworks)
- Database tests use SQLite to avoid requiring PostgreSQL for unit tests
- Test files follow Go convention: `*_test.go`