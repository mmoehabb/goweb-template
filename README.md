A template for developing full-stack web applications in Golang.

![overview](./overview.gif)

## Used Technologies

- [Golang](https://go.dev/)
- [Templ](https://templ.guide/)
- [Tailwind](https://tailwindcss.com/)
- [HTMX](https://htmx.org/)
- [Fiber](https://docs.gofiber.io/)
- [Postgres](https://github.com/jackc/pgx)

## Template Structure

### root files

#### main.go
As you probably know `main.go` is the starting point of your go application. In this template, the server is started and its endpoints declared explicitly in `main.go`, however it shall have direct access only to [hanlders](#handlers) which it uses in defining the endpoints (attaching fiber handlers).

#### cmd
Just an extra tool used as a shorthand for commands, as shown [below](#cmd-only-linux). You can safely delete it.

### handlers
This package (directory) includes all fiber callback functions, used in `main.go`, aggregated or grouped into different packages (directories). And for each sub-package there should exist two files: `types.go` and `validators.go`; the first defines related types to the group (i.e. User, Credentials...etc), whereas the latter defines different validate functions to be used in handlers while getting users inputs (requests payloads).

### db
This package exports functions related to the database, and more importantly it has sub-packages with a more specified functions for various (system) entities.

The exported functions are: Query, SeqQuery, Queries, and Disconnect. `Query` function establishes a connection then executes the supplied sql query (a string parameter), then returns the result as a slice of type `[]any`, and finally closes the connection. `SeqQuery` is just like `Query` however it doesn't close the connection eventually; so make sure to call `Disconnect` after a series (sequence) of `SeqQuery` calls. And lastly, `Queries` just takes a slice of strings parameter rather than a single string, executes each sql query, closes the connection, and then returns nil or error in case an error occured.

db package constitutes of several sub-packages each of which declares and defines various db functions, related to a specific entity, by using the above-mentioned four functions. For example, the `db/user` package exports `Add` and `Get` functions that can be used directly by [handlers](#handlers) to communicate with the database.

### public
Public assests live in the `./public` directory, which is served by the file server of fiber by using the method: `app.Static(...)` in `main.go`. You should put here all the pictures, videos, sound, scripts...etc, that shall be publicly served to all users with no restrictions. 

### pages & ui
These two packages contain only '.templ' files. As the name indicates, the first is for the application pages: templ components with '<head>' and '<body>' tag names. The latter, on the other hand, constitutes of several sub-packages for several ui categories, like: forms, components, layouts, mini-components...etc.

```
.
├── cmd
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── README.md
├── db
│   ├── db.go
│   └── user
│       └── queries.go
├── handlers
│   └── user
│       ├── login.go
│       ├── register.go
│       ├── types.go
│       └── validators.go
├── pages
│   └── index.templ
├── public
│   ├── globals.css
│   ├── tailwind.js
│   ├── util.js
│   └── ...
└── ui
    ├── components
    │   ├── Button.templ
    │   └── TextInput.templ
    ├── forms
    │   ├── login.templ
    │   └── register.templ
    └── layouts
        ├── footer.templ
        └── header.templ
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
    - specify user to "postgres": `$ su - postgres`
    - run the following command in order to find the coniguration file location: `$ psql -c "SHOW config_file"`
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

### ./cmd (Only Linux)

You may use the `./cmd` file placed in the root directory, as a shorthand for the above-mentioned commands:

```shell
$ chmod +x ./cmd
$ ./cmd dev # executes: "templ generate --watch --cmd 'air'"
```
