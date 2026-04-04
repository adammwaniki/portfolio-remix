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
	"strings"
	"time"

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
		"add":     func(a, b int) int { return a + b },
		"tagSlug": content.TagSlug,
		"join":    strings.Join,
		"lower":   strings.ToLower,
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
			"Page": "home", "Title": "Adam Mwaniki | Software Engineer",
			"Description":  "Adam Mwaniki | Software Engineer. I build software that's clear, maintainable and built to last.",
			"CanonicalURL": content.SiteURL + "/",
			"OGType":       "website",
			"Sections":     sections, "IsDetail": false, "IsDark": false,
		}},
	}

	for _, s := range sections {
		pages = append(pages, page{
			"/" + s.ID, map[string]any{
				"Page": s.ID, "Title": fmt.Sprintf("%s | Adam Mwaniki", s.Title),
				"Description":  s.Subtitle,
				"CanonicalURL": content.SiteURL + "/" + s.ID,
				"OGType":       "website",
				"Section":      s, "Sections": sections, "IsDetail": false,
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
					"Page": s.ID, "Title": fmt.Sprintf("%s | Adam Mwaniki", c.Title),
					"Description":  c.Description,
					"CanonicalURL": content.SiteURL + "/" + s.ID + "/" + c.ID,
					"OGType":       "article",
					"Section":      s, "Card": c, "Sections": sections,
					"IsDetail": true, "BackURL": "/" + s.ID,
					"IsDark": s.IsDark, "CardIndex": i + 1,
					"NextCard":     nextCard, "PrevCard": prevCard,
					"RelatedCards": content.RelatedCards(s.ID, c.ID, 3),
				},
			})
		}
	}

	// Tags index page
	pages = append(pages, page{
		"/tags", map[string]any{
			"Page": "tags", "Title": "Tags | Adam Mwaniki",
			"Description":  "Browse articles by topic across all sections.",
			"CanonicalURL": content.SiteURL + "/tags",
			"OGType":       "website",
			"AllTags":      content.AllTags(),
			"Sections":     sections, "IsDetail": false, "IsDark": false,
		},
	})

	// Individual tag pages
	for _, tag := range content.AllTags() {
		slug := content.TagSlug(tag)
		pages = append(pages, page{
			"/tags/" + slug, map[string]any{
				"Page": "tags", "Title": fmt.Sprintf("%s | Adam Mwaniki", tag),
				"Description":  fmt.Sprintf("Articles tagged with \"%s\" on Adam Mwaniki's portfolio.", tag),
				"CanonicalURL": content.SiteURL + "/tags/" + slug,
				"OGType":       "website",
				"Tag":          tag,
				"TagCards":     content.CardsByTag(tag),
				"Sections":     sections, "IsDetail": false, "IsDark": false,
			},
		})
	}

	// Contact page
	pages = append(pages, page{
		"/contact", map[string]any{
			"Page": "contact", "Title": "Contact | Adam Mwaniki",
			"Description":  "Get in touch with Adam Mwaniki | Software Engineer.",
			"CanonicalURL": content.SiteURL + "/contact",
			"OGType":       "website",
			"Sections":     sections, "IsDetail": false, "IsDark": false,
		},
	})

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

	// Copy static assets
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

	// Generate search index
	generateSearchIndex(sections)

	// Generate sitemap.xml
	generateSitemap()

	// Generate robots.txt
	generateRobots()

	// Generate RSS feed
	generateRSS(sections)

	// Generate worker
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

func generateSearchIndex(sections []content.Section) {
	var index []map[string]any
	for _, s := range sections {
		for _, c := range s.Cards {
			index = append(index, map[string]any{
				"title":       c.Title,
				"description": c.Description,
				"tags":        c.Tags,
				"url":         "/" + s.ID + "/" + c.ID,
				"section":     s.Title,
			})
		}
	}
	data, _ := json.Marshal(index)
	writeFile("dist/search-index.json", data)
}

func generateSitemap() {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	sections := content.Sections()
	now := time.Now().Format("2006-01-02")

	// Home
	buf.WriteString(fmt.Sprintf("  <url><loc>%s/</loc><lastmod>%s</lastmod><changefreq>weekly</changefreq><priority>1.0</priority></url>\n", content.SiteURL, now))

	// Section pages
	for _, s := range sections {
		buf.WriteString(fmt.Sprintf("  <url><loc>%s/%s</loc><lastmod>%s</lastmod><changefreq>weekly</changefreq><priority>0.8</priority></url>\n", content.SiteURL, s.ID, now))
		for _, c := range s.Cards {
			lastmod := c.Date
			if c.Updated != "" {
				lastmod = c.Updated
			}
			buf.WriteString(fmt.Sprintf("  <url><loc>%s/%s/%s</loc><lastmod>%s</lastmod><changefreq>monthly</changefreq><priority>0.6</priority></url>\n", content.SiteURL, s.ID, c.ID, lastmod))
		}
	}

	// Tags and contact
	buf.WriteString(fmt.Sprintf("  <url><loc>%s/tags</loc><lastmod>%s</lastmod><changefreq>weekly</changefreq><priority>0.5</priority></url>\n", content.SiteURL, now))
	buf.WriteString(fmt.Sprintf("  <url><loc>%s/contact</loc><lastmod>%s</lastmod><changefreq>monthly</changefreq><priority>0.5</priority></url>\n", content.SiteURL, now))

	buf.WriteString("</urlset>\n")
	writeFile("dist/sitemap.xml", buf.Bytes())
}

func generateRobots() {
	robots := fmt.Sprintf("User-agent: *\nAllow: /\n\nSitemap: %s/sitemap.xml\n", content.SiteURL)
	writeFile("dist/robots.txt", []byte(robots))
}

func generateRSS(sections []content.Section) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">` + "\n")
	buf.WriteString("<channel>\n")
	buf.WriteString("  <title>Adam Mwaniki</title>\n")
	buf.WriteString(fmt.Sprintf("  <link>%s</link>\n", content.SiteURL))
	buf.WriteString("  <description>Software Engineer. I build software that's clear, maintainable and built to last.</description>\n")
	buf.WriteString("  <language>en</language>\n")
	buf.WriteString(fmt.Sprintf("  <atom:link href=\"%s/feed.xml\" rel=\"self\" type=\"application/rss+xml\"/>\n", content.SiteURL))
	buf.WriteString(fmt.Sprintf("  <lastBuildDate>%s</lastBuildDate>\n", time.Now().Format(time.RFC1123Z)))

	// Collect all cards with dates, sort by date descending
	type entry struct {
		section content.Section
		card    content.Card
	}
	var entries []entry
	for _, s := range sections {
		for _, c := range s.Cards {
			entries = append(entries, entry{s, c})
		}
	}
	// Sort newest first
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].card.Date > entries[i].card.Date {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	for _, e := range entries {
		pubDate := e.card.Date
		t, err := time.Parse("2006-01-02", pubDate)
		if err == nil {
			pubDate = t.Format(time.RFC1123Z)
		}
		link := fmt.Sprintf("%s/%s/%s", content.SiteURL, e.section.ID, e.card.ID)
		buf.WriteString("  <item>\n")
		buf.WriteString(fmt.Sprintf("    <title>%s</title>\n", escapeXML(e.card.Title)))
		buf.WriteString(fmt.Sprintf("    <link>%s</link>\n", link))
		buf.WriteString(fmt.Sprintf("    <guid>%s</guid>\n", link))
		buf.WriteString(fmt.Sprintf("    <pubDate>%s</pubDate>\n", pubDate))
		buf.WriteString(fmt.Sprintf("    <description>%s</description>\n", escapeXML(e.card.Description)))
		buf.WriteString(fmt.Sprintf("    <category>%s</category>\n", escapeXML(e.section.Title)))
		buf.WriteString("  </item>\n")
	}

	buf.WriteString("</channel>\n</rss>\n")
	writeFile("dist/feed.xml", buf.Bytes())
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func generateWorker(titles map[string]string, fragments map[string]string) {
	titlesJSON, _ := json.Marshal(titles)
	fragmentsJSON, _ := json.Marshal(fragments)

	worker := fmt.Sprintf(`const TITLES = %s;
const FRAGMENTS = %s;

export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    // Static assets
    if (url.pathname.startsWith("/static/")) {
      return env.ASSETS.fetch(request);
    }

    // Sitemap, robots, RSS feed, search index
    if (url.pathname === "/sitemap.xml" || url.pathname === "/robots.txt" || url.pathname === "/feed.xml" || url.pathname === "/search-index.json") {
      return env.ASSETS.fetch(request);
    }

    // View counter API
    if (url.pathname.startsWith("/api/views/")) {
      return handleViews(request, env, url);
    }

    // HTMX partial requests
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

async function handleViews(request, env, url) {
  const path = url.pathname.replace("/api/views", "");
  const headers = new Headers({
    "Content-Type": "application/json",
    "Access-Control-Allow-Origin": "*",
  });

  if (!env.VIEWS) {
    return new Response(JSON.stringify({ views: 0 }), { headers });
  }

  try {
    if (request.method === "POST" || request.method === "GET") {
      const key = "views:" + path;
      const current = parseInt(await env.VIEWS.get(key) || "0", 10);
      if (request.method === "GET") {
        // Increment on GET for simplicity (each page view)
        const next = current + 1;
        await env.VIEWS.put(key, next.toString());
        return new Response(JSON.stringify({ views: next }), { headers });
      }
    }
  } catch (e) {
    return new Response(JSON.stringify({ views: 0 }), { headers });
  }

  return new Response(JSON.stringify({ views: 0 }), { headers });
}
`, string(titlesJSON), string(fragmentsJSON))

	writeFile("src/worker.js", []byte(worker))
}
