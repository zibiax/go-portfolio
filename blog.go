package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
)

var slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(title string) string {
	s := strings.ToLower(title)
	s = slugInvalidChars.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "post"
	}
	return s
}

func uniqueSlug(title string) (string, error) {
	base := slugify(title)
	slug := base
	for n := 2; ; n++ {
		exists, err := slugExists(slug)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, n)
	}
}

func renderMarkdown(source string) (template.HTML, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// Public handlers

type postView struct {
	Post
	ContentHTML template.HTML
}

func handleBlogList(w http.ResponseWriter, r *http.Request) {
	posts, err := listPublishedPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	views := make([]postView, 0, len(posts))
	for _, p := range posts {
		html, err := renderMarkdown(p.ContentMarkdown)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		views = append(views, postView{Post: p, ContentHTML: html})
	}

	tmpl, err := template.ParseFiles("templates/blog_list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]any{"Posts": views})
}

func handleBlogPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := getPostBySlug(slug)
	if errors.Is(err, sql.ErrNoRows) || (post != nil && !post.Published) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html, err := renderMarkdown(post.ContentMarkdown)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/blog_post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]any{"Post": post, "ContentHTML": html})
}

// Admin handlers

func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	posts, err := listAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin_dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]any{"Posts": posts})
}

func handleAdminPostNewForm(w http.ResponseWriter, r *http.Request) {
	renderPostForm(w, nil, "")
}

func handleAdminPostNewSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	published := r.FormValue("published") == "on"

	if title == "" || content == "" {
		renderPostForm(w, nil, "Title and content are required")
		return
	}

	slug, err := uniqueSlug(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := createPost(title, slug, content, published); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleAdminPostEditForm(w http.ResponseWriter, r *http.Request) {
	post, err := loadPostFromPath(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderPostForm(w, post, "")
}

func handleAdminPostEditSubmit(w http.ResponseWriter, r *http.Request) {
	post, err := loadPostFromPath(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	published := r.FormValue("published") == "on"

	if title == "" || content == "" {
		renderPostForm(w, post, "Title and content are required")
		return
	}

	if err := updatePost(post.ID, title, content, published); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleAdminPostDelete(w http.ResponseWriter, r *http.Request) {
	post, err := loadPostFromPath(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := deletePost(post.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func loadPostFromPath(r *http.Request) (*Post, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return getPostByID(id)
}

func renderPostForm(w http.ResponseWriter, post *Post, errMsg string) {
	tmpl, err := template.ParseFiles("templates/admin_post_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]any{"Post": post, "Error": errMsg})
}
