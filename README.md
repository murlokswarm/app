<h1 align="center">
    <a href="https://github.com/maxence-charriere/app">
        <img alt="app"  width="150" height="150" src="https://storage.googleapis.com/murlok-github/icon-192.png">
    </a>
</h1>

<p align="center">
	<a href="https://circleci.com/gh/maxence-charriere/app"><img src="https://circleci.com/gh/maxence-charriere/app.svg?style=svg" alt="Circle CI Go build"></a>
    <a href="https://goreportcard.com/report/github.com/maxence-charriere/app"><img src="https://goreportcard.com/badge/github.com/maxence-charriere/app" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/maxence-charriere/app/pkg/app"><img src="https://godoc.org/github.com/maxence-charriere/app/pkg/app?status.svg" alt="GoDoc"></a>
</p>

**app** is a package to build [progressive web apps (PWA)](https://developers.google.com/web/progressive-web-apps/) with [Go programming language](https://golang.org) and [WebAssembly](https://webassembly.org).

It uses a [declarative syntax](#declarative-syntax) that allows creating and dealing with HTML elements only by using Go, and without writing any HTML markup.

The package also provides an [http.handler](#http-handler) ready to serve all the required resources to run Go-based progressive web apps.

## Install

**app** requires [Go 1.13](https://golang.org/doc/go1.13) or newer.

```sh
go get -u -v github.com/maxence-charriere/app/pkg/app
```

## Declarative syntax

**app** uses a declarative syntax so you can write component-based UI elements just by using the Go programming language.

```go
package main

import "github.com/maxence-charriere/app/pkg/app"

type hello struct {
    app.Compo
    name string
}

func (h *hello) Render() app.UI {
    return app.Div().Body(
        app.Main().Body(
            app.H1().Body(
                app.Text("Hello, "),
                app.If(h.name != "",
                    app.Text(h.name),
                ).Else(
                    app.Text("World"),
                ),
            ),
            app.Input().
                Value(h.name).
                Placeholder("What is your name?").
                AutoFocus(true).
                OnChange(h.OnInputChange),
        ),
    )
}

func (h *hello) OnInputChange(src app.Value, e app.Event) {
	h.name = src.Get("value").String()
	h.Update()
}

func main() {
	app.Route("/", &hello{})
	app.Run()
}

```

Then you can build your user interface by using the Go build tool:

```sh
GOOS=js GOARCH=wasm go build -o app.wasm
```

## HTTP handler

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/maxence-charriere/app/pkg/app"
)

func main() {
    fmt.Println("starting local server")

    h := &app.Handler{
        Title:  "Hello Demo",
        Author: "Maxence Charriere",
    }

    if err := http.ListenAndServe(":7777", h); err != nil {
        panic(err)
    }
}
```

## Works on mainstream browsers

|         | Chrome | Edge | Firefox | Opera | Safari |
| :------ | :----: | :--: | :-----: | :---: | :----: |
| Desktop |   ✔    | ✔\*  |    ✔    |   ✔   |   ✔    |
| Mobile  |   ✔    |  ✔   |    ✔    |   ✔   |   ✔    |

\*_only Chromium based [Edge](https://www.microsoft.com/edge)_

## Live demo

- [Luck](https://luck.murlok.io)
- [Hello](https://demo.murlok.io)
- [City](https://demo.murlok.io/city)
