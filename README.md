A template for developing full-stack web applications in Golang.

## Used Technologies

- [Golang](https://go.dev/)
- [Templ](https://templ.guide/)
- [Tailwind](https://tailwindcss.com/)
- [HTMX](https://htmx.org/)
- [Fiber](https://docs.gofiber.io/)

## Usage

Download the source code or just clone this repository and delete .git directory:

```shell
$ git clone https://github.com/mmoehabb/goweb-template
$ rm -rf .git
$ git init # optional
```

> Make sure you have installed [go](https://go.dev/doc/install) and [templ](https://templ.guide/quick-start/installation):

Execute the following command on the root directory:

```shell
$ templ generate --cmd "go run ."
```

This shall generate `.go` files from `.templ` files, and run the server afterwards.
If everything went right, you should be able to see the template live on [http://localhost:3000](http://localhost:3000)

You can also enable live reload with the command:

```
$ templ generate --watch --cmd "go run ."
```
