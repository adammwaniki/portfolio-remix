package main

import (
    "log"
    "net/http"

	"github.com/adammwaniki/portfolio-remix/views"
)

func main() {
    mux := http.NewServeMux()

    // Route for homepage
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering Home Page")
		views.Home().Render(r.Context(), w)
    })
	// Route for projects page
    mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Rendering Projects Page")
		views.Page(true).Render(r.Context(), w)
    })

    // Serve Tailwind CSS and any other static assets
    mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	mux.Handle("/favicon.ico", http.FileServer(http.Dir("public")))


    log.Println("Listening on http://localhost:8080")
    http.ListenAndServe(":8080", mux)
}
