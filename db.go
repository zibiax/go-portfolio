package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	content_markdown TEXT NOT NULL,
	published INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
	token_hash TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL,
	expires_at DATETIME NOT NULL
);
`

func initDB(path string) error {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	if _, err := conn.Exec(schema); err != nil {
		conn.Close()
		return err
	}

	db = conn
	return nil
}

type Post struct {
	ID              int64
	Title           string
	Slug            string
	ContentMarkdown string
	Published       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func createPost(title, slug, contentMarkdown string, published bool) (int64, error) {
	now := time.Now().UTC()
	res, err := db.Exec(
		`INSERT INTO posts (title, slug, content_markdown, published, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		title, slug, contentMarkdown, published, now, now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updatePost(id int64, title, contentMarkdown string, published bool) error {
	_, err := db.Exec(
		`UPDATE posts SET title = ?, content_markdown = ?, published = ?, updated_at = ? WHERE id = ?`,
		title, contentMarkdown, published, time.Now().UTC(), id,
	)
	return err
}

func deletePost(id int64) error {
	_, err := db.Exec(`DELETE FROM posts WHERE id = ?`, id)
	return err
}

func slugExists(slug string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM posts WHERE slug = ?)`, slug).Scan(&exists)
	return exists, err
}

func getPostByID(id int64) (*Post, error) {
	row := db.QueryRow(
		`SELECT id, title, slug, content_markdown, published, created_at, updated_at FROM posts WHERE id = ?`, id,
	)
	return scanPost(row)
}

func getPostBySlug(slug string) (*Post, error) {
	row := db.QueryRow(
		`SELECT id, title, slug, content_markdown, published, created_at, updated_at FROM posts WHERE slug = ?`, slug,
	)
	return scanPost(row)
}

func scanPost(row *sql.Row) (*Post, error) {
	var p Post
	err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.ContentMarkdown, &p.Published, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func listAllPosts() ([]Post, error) {
	rows, err := db.Query(
		`SELECT id, title, slug, content_markdown, published, created_at, updated_at FROM posts ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectPosts(rows)
}

func listPublishedPosts() ([]Post, error) {
	rows, err := db.Query(
		`SELECT id, title, slug, content_markdown, published, created_at, updated_at FROM posts WHERE published = 1 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectPosts(rows)
}

func collectPosts(rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.ContentMarkdown, &p.Published, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
