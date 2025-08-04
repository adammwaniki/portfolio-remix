package main

import (
    "os"
    "log"
    "github.com/adammwaniki/portfolio-remix/views"

    "github.com/a-h/templ"
)

func main() {
    f, err := os.Create("dist/index.html")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    err = templ.Execute(f, views.Page())
    if err != nil {
        log.Fatal(err)
    }
}
