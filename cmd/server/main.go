package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stefannegele/website/internal/db"
	"github.com/stefannegele/website/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "blog.db"
	}

	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("opening db: %v", err)
	}
	defer database.Close()

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2 January 2006")
		},
		"sourceLabel": func(s string) string {
			switch s {
			case "innoq":
				return "INNOQ"
			case "external":
				return "External"
			default:
				return "Article"
			}
		},
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob(filepath.Join("templates", "*.html")))
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("templates", "admin", "*.html")))

	h := handlers.New(database, tmpl)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public routes
	r.Get("/", h.Index)
	r.Get("/posts/{slug}", h.PostDetail)
	r.Get("/impressum", h.Impressum)

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Get("/", h.AdminIndex)
		r.Get("/new", h.AdminNewPost)
		r.Post("/new", h.AdminCreatePost)
		r.Get("/{id}/edit", h.AdminEditPost)
		r.Post("/{id}/edit", h.AdminUpdatePost)
		r.Post("/{id}/delete", h.AdminDeletePost)
		r.Post("/sync-rss", h.AdminSyncRSS)
	})

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
