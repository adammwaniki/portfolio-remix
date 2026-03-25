package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adammwaniki/portfolio-remix/internal/content"
)

var tmpl *template.Template

func main() {
	// Parse all templates
	var err error
	tmpl, err = template.New("").Funcs(template.FuncMap{
		"raw": func(s string) template.HTML { return template.HTML(s) },
		"dict": func(values ...any) map[string]any {
			m := make(map[string]any)
			for i := 0; i < len(values)-1; i += 2 {
				key, _ := values[i].(string)
				m[key] = values[i+1]
			}
			return m
		},
	}).ParseGlob("views/*.html")
	if err != nil {
		log.Fatal("parsing views: ", err)
	}

	partials, err := tmpl.Clone()
	if err != nil {
		log.Fatal("cloning templates: ", err)
	}
	tmpl, err = partials.ParseGlob("views/partials/*.html")
	if err != nil {
		log.Fatal("parsing partials: ", err)
	}

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/technical-notes", handleSection)
	http.HandleFunc("/projects", handleSection)
	http.HandleFunc("/musings", handleSection)
	http.HandleFunc("/the-bullshitters", handleSection)

	// Detail routes: /section-id/card-id
	http.HandleFunc("/technical-notes/", handleDetail)
	http.HandleFunc("/projects/", handleDetail)
	http.HandleFunc("/musings/", handleDetail)
	http.HandleFunc("/the-bullshitters/", handleDetail)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// isHTMX checks if the request was made by HTMX.
func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// renderPage renders a full page or just the fragment if HTMX.
func renderPage(w http.ResponseWriter, r *http.Request, fragmentName string, data any) {
	if isHTMX(r) && r.Header.Get("HX-History-Restore-Request") != "true" {
		// Send title via header so JS can update document.title
		if d, ok := data.(map[string]any); ok {
			if title, ok := d["Title"].(string); ok {
				w.Header().Set("HX-Title", title)
			}
		}
		// Return the header + main + footer fragment for swap
		if err := tmpl.ExecuteTemplate(w, "page-content", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	// Full page: wrap in layout
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]any{
		"Page":     "home",
		"Title":    "Adam Mwaniki | Software Engineer",
		"Sections": content.Sections(),
		"IsDetail": false,
	}
	renderPage(w, r, "home.html", data)
}

func handleSection(w http.ResponseWriter, r *http.Request) {
	sectionID := strings.TrimPrefix(r.URL.Path, "/")

	section, ok := content.SectionByID(sectionID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := map[string]any{
		"Page":     sectionID,
		"Title":    fmt.Sprintf("%s | Adam Mwaniki", section.Title),
		"Section":  section,
		"Sections": content.Sections(),
		"IsDetail": false,
	}
	renderPage(w, r, "section.html", data)
}

func handleDetail(w http.ResponseWriter, r *http.Request) {
	// Parse /section-id/card-id
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		// Redirect to section if no card ID
		http.Redirect(w, r, "/"+parts[0], http.StatusFound)
		return
	}

	sectionID, cardID := parts[0], parts[1]
	section, card, ok := content.CardByID(sectionID, cardID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := map[string]any{
		"Page":      sectionID,
		"Title":     fmt.Sprintf("%s | Adam Mwaniki", card.Title),
		"Section":   section,
		"Card":      card,
		"Sections":  content.Sections(),
		"IsDetail":  true,
		"BackURL":   "/" + sectionID,
	}
	renderPage(w, r, "detail.html", data)
}