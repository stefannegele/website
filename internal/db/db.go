package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

type Post struct {
	ID          int64
	Title       string
	Slug        string
	Summary     string
	Content     string
	URL         string
	Source      string // "innoq", "external", "internal"
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	db := &DB{sqlDB}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) migrate() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			title       TEXT NOT NULL,
			slug        TEXT UNIQUE NOT NULL,
			summary     TEXT,
			content     TEXT,
			url         TEXT,
			source      TEXT NOT NULL CHECK(source IN ('innoq','external','internal')),
			published_at DATETIME NOT NULL,
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published_at DESC);
		CREATE INDEX IF NOT EXISTS idx_posts_source ON posts(source);
	`)
	return err
}

func (db *DB) AllPosts() ([]Post, error) {
	rows, err := db.Query(`
		SELECT id, title, slug, summary, content, url, source, published_at, created_at, updated_at
		FROM posts ORDER BY published_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

func (db *DB) PostBySlug(slug string) (*Post, error) {
	row := db.QueryRow(`
		SELECT id, title, slug, summary, content, url, source, published_at, created_at, updated_at
		FROM posts WHERE slug = ?
	`, slug)
	p := &Post{}
	err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.Summary, &p.Content, &p.URL, &p.Source, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (db *DB) UpsertInnoqPost(p *Post) error {
	_, err := db.Exec(`
		INSERT INTO posts (title, slug, summary, content, url, source, published_at)
		VALUES (?, ?, ?, ?, ?, 'innoq', ?)
		ON CONFLICT(slug) DO UPDATE SET
			title=excluded.title,
			summary=excluded.summary,
			content=excluded.content,
			url=excluded.url,
			published_at=excluded.published_at,
			updated_at=CURRENT_TIMESTAMP
	`, p.Title, p.Slug, p.Summary, p.Content, p.URL, p.PublishedAt)
	return err
}

func (db *DB) CreatePost(p *Post) error {
	_, err := db.Exec(`
		INSERT INTO posts (title, slug, summary, content, url, source, published_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, p.Title, p.Slug, p.Summary, p.Content, p.URL, p.Source, p.PublishedAt)
	return err
}

func (db *DB) UpdatePost(p *Post) error {
	_, err := db.Exec(`
		UPDATE posts SET title=?, slug=?, summary=?, content=?, url=?, published_at=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?
	`, p.Title, p.Slug, p.Summary, p.Content, p.URL, p.PublishedAt, p.ID)
	return err
}

func (db *DB) DeletePost(id int64) error {
	_, err := db.Exec(`DELETE FROM posts WHERE id=?`, id)
	return err
}

func (db *DB) PostByID(id int64) (*Post, error) {
	row := db.QueryRow(`
		SELECT id, title, slug, summary, content, url, source, published_at, created_at, updated_at
		FROM posts WHERE id = ?
	`, id)
	p := &Post{}
	err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.Summary, &p.Content, &p.URL, &p.Source, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func scanPosts(rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Summary, &p.Content, &p.URL, &p.Source, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
