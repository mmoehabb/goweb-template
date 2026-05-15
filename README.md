A template for developing full-stack web applications in Golang.

![overview](./overview.gif)

## Used Technologies

- [Golang](https://go.dev/)
- [Templ](https://templ.guide/)
- [TemplUI](https://templui.io/)
- [Tailwind](https://tailwindcss.com/)
- [HTMX](https://htmx.org/)
- [Fiber](https://docs.gofiber.io/)
- [Postgres](https://github.com/jackc/pgx)
- [GORM](https://gorm.io/) (Database ORM)
- [Goose](https://github.com/pressly/goose) (SQL Migrations)

## Template Structure

### root files

#### main.go
As you probably know `main.go` is the starting point of your go application. In this template, it merely initializes a fiber app, adds some middlewares, and defines a couple of endpoints.

#### luci.config.toml
Just an extra tool used as a shorthand for commands, as shown [below](#luci-cli). You can just delete it, if you don't like it!

### handlers
This package (directory) includes all fiber callback functions, used in `main.go`, aggregated or grouped into different packages (directories). And for each sub-package there should exist two files: `types.go` and `validators.go`; the first defines related types to the group (i.e. User, Credentials...etc), whereas the latter defines different validate functions to be used in handlers while getting users inputs (requests payloads).

### db
This package exports functions related to the database using [GORM](https://gorm.io/). It provides a global `Connection` variable for database operations and sub-packages with entity-specific functions.

The main exports are:
- `Init()` - initializes the GORM database connection
- `RunMigrations(models...)` - auto-migrates model schemas
- `Connection` - global `*gorm.DB` instance for direct queries

For example, the `db/users` package exports `Add` and `Get` functions that can be used directly by [handlers](#handlers) to communicate with the database.

### public
Public assests live in the `./public` directory, which is served by the file server of fiber by using the method: `app.Static(...)` in `main.go`. You should put here all the pictures, videos, sound, scripts...etc, that shall be publicly served to all users with no restrictions. 

### pages & ui
These two packages contain only '.templ' files. As the name indicates, the first is for the application pages: templ components with '<head>' and '<body>' tag names. The latter, on the other hand, constitutes of several sub-packages for several ui categories, like: forms, components, layouts, mini-components...etc.

### ancillaries
This package shall contain all logic that's shared between other packages and sub-packages.

### constants

All constant values shall be defined in this package. For example, your `.env` file values are represented in this package as a global go struct that can be imported from anywhere else.

### Directory Tree

```
.
├── luci.config.toml
├── go.mod
├── go.sum
├── .templui.json
├── LICENSE
├── main.go
├── README.md
├── ancillaries
│   └── errors.go
├── components
│   └── button
│       └── button.templ
├── constants
│   └── config.go
├── db
│   ├── db.go
│   ├── goose_migrations/
│   │   └── 001_create_users.sql
│   └── users/
│       ├── model.go
│       └── queries.go
├── handlers
│   └── user
│       ├── login.go
│       ├── register.go
│       ├── types.go
│       └── validators.go
├── pages
│   └── index.templ
├── public
│   ├── globals.css
│   ├── tailwind.js
│   ├── util.js
│   └── ...
├── ui
│   ├── components
│   │   ├── Button.templ
│   │   └── TextInput.templ
│   ├── forms
│   │   ├── login.templ
│   │   └── register.templ
│   └── layouts
│       ├── footer.templ
│       └── header.templ
└── utils
    └── templui.go
```

## Usage

Download the source code or just clone this repository and delete .git directory:

```shell
$ git clone https://github.com/mmoehabb/goweb-template
$ rm -rf .git
$ git init # optional
```

> Make sure you have installed [go](https://go.dev/doc/install) and [templ](https://templ.guide/quick-start/installation):

Install the dependencies with; execute the following command on the root directory:

```shell
$ go install
```

> You may use `luci install`, as mentioned below in "Luci CLI" section, to install both the packages and the tools (templ and air binaries) all at once.

Then, write the following command to compile templ files and run the server afterwards:

```shell
$ templ generate --cmd "go run ."
```

If everything went right, you should be able to see the template live on [http://localhost:3000](http://localhost:3000)

You can also enable live reload with the command:

```shell
$ templ generate --watch --cmd "go run ."
```

However this will watch only templ files, you may wanna reload the server when go files are modified as well.
For this sake `.air.toml` file (as you may have noticed) is in the root directory; make sure to install [air](https://github.com/air-verse/air) then execute the previous command with `air` instead of `go run .`.

```shell
$ templ generate --watch --cmd "air"
```

### Play the Game

In order to see all the template functionalities in action, you have to make sure that a postgresql service is running on your machine with a 'postgres' database created, or any (postgres) database you already have but make sure to modify the database connection config in `db/db.go` file accordingly:

> If you just wanna PLAY THE GAME without all this crap; just call `start()` function in the developer tools.

```Golang
// You will find this line in `connect` function in `db.go`.
conn, err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
```

If you haven't established a postgresql server before, you may find the following steps helpful:
1. Download & install postgres from here: [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
2. Modify `pg_hba.conf` to enable md5 remote access:
    - log into the terminal with "postgres" user: `$ su - postgres`
    - run the following command in order to find the configuration file location: `$ psql -c "SHOW config_file"`
    - open the file `pg_hba.conf` located at the same directory of `postgresql.conf`, then add the lines, shown below step 3, to the end of it:
3. Start the service:
    ```shell
    $ service postgresql start
    ```

```
# Add these lines to the end of pg_hba.conf
# This means that remote access is allowed using IP v4 and IP v6 to all databases and all users using the "md5" authentication protocol
host    all             all              127.0.0.1/0                       md5
host    all             all              ::/0                            md5
```

And finally run the application, register, login, and have fun:

```shell
$ go run .
```

### Database Migrations

This template supports **two migration strategies** that run automatically on startup:

#### 1. GORM AutoMigrate (Recommended for prototyping)
For quick schema changes during development. Define models with GORM tags:

```go
package users

import "gorm.io/gorm"

type DataModel struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;not null"`
    Password string `gorm:"not null"`
}

// TableName overrides the default table name (data_model; which is generated
// from the type name) to "users"
func (DataModel) TableName() string {
	return "users"
}
```

Add models to `db.RunMigrations()` in `main.go`:

```go
if err := db.RunMigrations(users.DataModel{}); err != nil {
    log.Fatalf("Failed to run GORM migrations: %v", err)
}
```

#### 2. Goose SQL Migrations (Recommended for production)
For version-controlled SQL migrations with rollback support. SQL files go in `db/goose_migrations/`:

```sql
-- db/goose_migrations/001_create_users.sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(45) UNIQUE NOT NULL,
    password VARCHAR(45) NOT NULL,
    created_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
```

Goose migrations also run automatically on startup via `db.RunGooseMigrations()`.

#### When to use which?
- **GORM AutoMigrate**: Fast prototyping, simple models, don't need rollbacks
- **Goose**: Production apps, complex SQL, need version control, require rollbacks

**Note:** Goose CLI commands use the `DATABASE_URL` from your `.env` file. Make sure to source it before running migrate commands:

```shell
$ source .env
$ DATABASE_URL=$DATABASE_URL luci migrate up
```

### Luci CLI

You may use the [luci](https://github.com/mmoehabb/luci) CLI tool, as a shorthand for the above-mentioned commands:

```shell
$ go install github.com/mmoehabb/luci@latest
$ luci dev # executes: "templ generate --watch --cmd 'air'"
```
