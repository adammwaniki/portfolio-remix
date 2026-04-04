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
		"add": func(a, b int) int { return a + b },
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

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/technical-notes", handleSection)
	http.HandleFunc("/projects", handleSection)
	http.HandleFunc("/musings", handleSection)
	http.HandleFunc("/the-bullshitters", handleSection)

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

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func renderPage(w http.ResponseWriter, r *http.Request, data any) {
	if isHTMX(r) && r.Header.Get("HX-History-Restore-Request") != "true" {
		if d, ok := data.(map[string]any); ok {
			if title, ok := d["Title"].(string); ok {
				w.Header().Set("HX-Title", title)
			}
		}
		if err := tmpl.ExecuteTemplate(w, "page-content", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
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
		"Title":    "Adam Mwaniki \u2014 Software Engineer",
		"Sections": content.Sections(),
		"IsDetail": false,
		"IsDark":   false,
	}
	renderPage(w, r, data)
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
		"Title":    fmt.Sprintf("%s \u2014 Adam Mwaniki", section.Title),
		"Section":  section,
		"Sections": content.Sections(),
		"IsDetail": false,
		"IsDark":   section.IsDark,
	}
	renderPage(w, r, data)
}

func handleDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		http.Redirect(w, r, "/"+parts[0], http.StatusFound)
		return
	}

	sectionID, cardID := parts[0], parts[1]
	section, card, ok := content.CardByID(sectionID, cardID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Find card index and next/prev for navigation
	cardIndex := 0
	var nextCard, prevCard map[string]any
	for i, c := range section.Cards {
		if c.ID == cardID {
			cardIndex = i
			if i > 0 {
				prevCard = map[string]any{
					"Title": section.Cards[i-1].Title,
					"URL":   "/" + sectionID + "/" + section.Cards[i-1].ID,
				}
			}
			if i < len(section.Cards)-1 {
				nextCard = map[string]any{
					"Title": section.Cards[i+1].Title,
					"URL":   "/" + sectionID + "/" + section.Cards[i+1].ID,
				}
			}
			break
		}
	}

	data := map[string]any{
		"Page":      sectionID,
		"Title":     fmt.Sprintf("%s \u2014 Adam Mwaniki", card.Title),
		"Section":   section,
		"Card":      card,
		"CardIndex": cardIndex + 1,
		"NextCard":  nextCard,
		"PrevCard":  prevCard,
		"Sections":  content.Sections(),
		"IsDetail":  true,
		"IsDark":    section.IsDark,
		"BackURL":   "/" + sectionID,
	}
	renderPage(w, r, data)
}
