package main

import (
    "context"
    "log"
    "os"
    "path/filepath"

    "github.com/adammwaniki/portfolio-remix/views"
)

func worker() {
    // Ensure the dist directory exists
    err := os.MkdirAll("dist", 0755)
    if err != nil {
        log.Fatal(err)
    }

    f, err := os.Create(filepath.Join("dist", "index.html"))
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    page := views.Page() // call the function to get templ.Component

    err = page.Render(context.Background(), f)
    if err != nil {
        log.Fatal(err)
    }
}
