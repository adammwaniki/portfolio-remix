package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/adammwaniki/portfolio-remix/internal/content"
)

func main() {
	funcMap := template.FuncMap{
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
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("views/*.html")
	if err != nil {
		log.Fatal("parsing views: ", err)
	}
	tmpl, err = tmpl.ParseGlob("views/partials/*.html")
	if err != nil {
		log.Fatal("parsing partials: ", err)
	}

	os.RemoveAll("dist")

	sections := content.Sections()

	type page struct {
		path string
		data map[string]any
	}

	pages := []page{
		{"/", map[string]any{
			"Page": "home", "Title": "Adam Mwaniki \u2014 Software Engineer",
			"Sections": sections, "IsDetail": false, "IsDark": false,
		}},
	}

	for _, s := range sections {
		pages = append(pages, page{
			"/" + s.ID, map[string]any{
				"Page": s.ID, "Title": fmt.Sprintf("%s \u2014 Adam Mwaniki", s.Title),
				"Section": s, "Sections": sections, "IsDetail": false,
				"IsDark": s.IsDark,
			},
		})
		for i, c := range s.Cards {
			var nextCard, prevCard map[string]any
			if i > 0 {
				prevCard = map[string]any{
					"Title": s.Cards[i-1].Title,
					"URL":   "/" + s.ID + "/" + s.Cards[i-1].ID,
				}
			}
			if i < len(s.Cards)-1 {
				nextCard = map[string]any{
					"Title": s.Cards[i+1].Title,
					"URL":   "/" + s.ID + "/" + s.Cards[i+1].ID,
				}
			}
			pages = append(pages, page{
				"/" + s.ID + "/" + c.ID, map[string]any{
					"Page": s.ID, "Title": fmt.Sprintf("%s \u2014 Adam Mwaniki", c.Title),
					"Section": s, "Card": c, "Sections": sections,
					"IsDetail": true, "BackURL": "/" + s.ID,
					"IsDark": s.IsDark, "CardIndex": i + 1,
					"NextCard": nextCard, "PrevCard": prevCard,
				},
			})
		}
	}

	titles := make(map[string]string)
	fragments := make(map[string]string)

	for _, p := range pages {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "layout.html", p.data); err != nil {
			log.Fatalf("rendering %s: %v", p.path, err)
		}
		writeFile(filepath.Join("dist", urlToFilePath(p.path)), buf.Bytes())

		buf.Reset()
		if err := tmpl.ExecuteTemplate(&buf, "page-content", p.data); err != nil {
			log.Fatalf("rendering fragment %s: %v", p.path, err)
		}
		fragments[p.path] = buf.String()
		titles[p.path] = p.data["Title"].(string)
	}

	filepath.WalkDir("static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		writeFile(filepath.Join("dist", path), data)
		return nil
	})

	generateWorker(titles, fragments)

	log.Printf("Built %d pages to dist/", len(pages))
}

func urlToFilePath(urlPath string) string {
	if urlPath == "/" {
		return "index.html"
	}
	return urlPath[1:] + "/index.html"
}

func writeFile(path string, data []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
}

func generateWorker(titles map[string]string, fragments map[string]string) {
	titlesJSON, _ := json.Marshal(titles)
	fragmentsJSON, _ := json.Marshal(fragments)

	worker := fmt.Sprintf(`const TITLES = %s;
const FRAGMENTS = %s;

export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname.startsWith("/static/")) {
      return env.ASSETS.fetch(request);
    }

    const isHtmx = request.headers.get("HX-Request") === "true";
    const isHistoryRestore = request.headers.get("HX-History-Restore-Request") === "true";

    if (isHtmx && !isHistoryRestore) {
      let path = url.pathname;
      if (path !== "/" && path.endsWith("/")) {
        path = path.slice(0, -1);
      }

      const fragment = FRAGMENTS[path];
      if (fragment) {
        const headers = new Headers();
        headers.set("Content-Type", "text/html; charset=utf-8");
        const title = TITLES[path];
        if (title) headers.set("HX-Title", title);
        return new Response(fragment, { status: 200, headers });
      }
    }

    return env.ASSETS.fetch(request);
  },
};
`, string(titlesJSON), string(fragmentsJSON))

	writeFile("src/worker.js", []byte(worker))
}
