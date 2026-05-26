package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stefannegele/website/internal/db"
	"github.com/stefannegele/website/internal/rss"
)

type Handler struct {
	DB        *db.DB
	Templates *template.Template
}

func New(database *db.DB, tmpl *template.Template) *Handler {
	return &Handler{DB: database, Templates: tmpl}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := h.DB.AllPosts()
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}
	h.render(w, "index.html", map[string]any{"Posts": posts})
}

func (h *Handler) PostDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := h.DB.PostBySlug(slug)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}
	h.render(w, "post.html", map[string]any{"Post": post})
}

func (h *Handler) Impressum(w http.ResponseWriter, r *http.Request) {
	h.render(w, "impressum.html", nil)
}

// Admin handlers
func (h *Handler) AdminIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := h.DB.AllPosts()
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}
	synced := r.URL.Query().Get("synced")
	h.render(w, "admin-index.html", map[string]any{"Posts": posts, "Synced": synced})
}

func (h *Handler) AdminNewPost(w http.ResponseWriter, r *http.Request) {
	h.render(w, "admin-edit.html", map[string]any{
		"Post":   &db.Post{PublishedAt: time.Now()},
		"IsNew":  true,
	})
}

func (h *Handler) AdminCreatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	post := postFromForm(r)
	post.Source = r.FormValue("source")
	if post.Source == "" {
		post.Source = "internal"
	}
	if err := h.DB.CreatePost(post); err != nil {
		h.render(w, "admin-edit.html", map[string]any{
			"Post":  post,
			"IsNew": true,
			"Error": err.Error(),
		})
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) AdminEditPost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	post, err := h.DB.PostByID(id)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}
	h.render(w, "admin-edit.html", map[string]any{"Post": post, "IsNew": false})
}

func (h *Handler) AdminUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	r.ParseForm()
	post := postFromForm(r)
	post.ID = id
	if err := h.DB.UpdatePost(post); err != nil {
		h.render(w, "admin-edit.html", map[string]any{
			"Post":  post,
			"IsNew": false,
			"Error": err.Error(),
		})
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) AdminDeletePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	h.DB.DeletePost(id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) AdminSyncRSS(w http.ResponseWriter, r *http.Request) {
	posts, err := rss.FetchInnoqPosts()
	if err != nil {
		http.Error(w, fmt.Sprintf("RSS sync failed: %v", err), 500)
		return
	}
	count := 0
	for i := range posts {
		if err := h.DB.UpsertInnoqPost(&posts[i]); err == nil {
			count++
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/admin?synced=%d", count), http.StatusSeeOther)
}

func (h *Handler) render(w http.ResponseWriter, tmpl string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Templates.ExecuteTemplate(w, tmpl, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func postFromForm(r *http.Request) *db.Post {
	pubStr := r.FormValue("published_at")
	published, err := time.Parse("2006-01-02", pubStr)
	if err != nil {
		published = time.Now()
	}
	return &db.Post{
		Title:       r.FormValue("title"),
		Slug:        slugify(r.FormValue("title")),
		Summary:     r.FormValue("summary"),
		Content:     r.FormValue("content"),
		URL:         r.FormValue("url"),
		PublishedAt: published,
	}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}
