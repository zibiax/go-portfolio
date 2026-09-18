package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/v35/github"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type Project struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Language    string    `json:"language"`
	Topics      []string  `json:"topics"`
	Stars       int       `json:"stars"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// projectsPerPage is how many GitHub repositories are fetched per request.
// The front page loads one page at a time behind a "Load more" button.
const projectsPerPage = 15

type ProjectsPage struct {
	Projects []Project `json:"projects"`
	Page     int       `json:"page"`
	HasMore  bool      `json:"hasMore"`
}

func getProjectData(repo *github.Repository) Project {
	return Project{
		Name:        repo.GetName(),
		Description: repo.GetDescription(),
		URL:         repo.GetHTMLURL(),
		Language:    repo.GetLanguage(),
		Topics:      repo.Topics,
		Stars:       repo.GetStargazersCount(),
		LastUpdated: repo.GetPushedAt().Time,
	}
}

func main() {
	genHash := flag.String("genhash", "", "print a bcrypt hash of the given password (for ADMIN_PASSWORD_HASH) and exit")
	flag.Parse()

	if *genHash != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*genHash), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(hash))
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("Warning: GITHUB_TOKEN is not set")
	}

	if os.Getenv("ADMIN_USERNAME") == "" || os.Getenv("ADMIN_PASSWORD_HASH") == "" {
		log.Println("Warning: ADMIN_USERNAME / ADMIN_PASSWORD_HASH not set, admin login will always fail")
	}

	dbPath := os.Getenv("BLOG_DB_PATH")
	if dbPath == "" {
		dbPath = "data/blog.db"
	}
	if err := initDB(dbPath); err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handleIndex)
	mux.HandleFunc("GET /projects", handleProjects)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /blog", handleBlogList)
	mux.HandleFunc("GET /blog/{slug}", handleBlogPost)

	mux.HandleFunc("GET /login", handleLoginPage)
	mux.HandleFunc("POST /login", handleLoginSubmit)
	mux.HandleFunc("POST /logout", handleLogout)

	mux.HandleFunc("GET /admin", requireAuth(handleAdminDashboard))
	mux.HandleFunc("GET /admin/posts/new", requireAuth(handleAdminPostNewForm))
	mux.HandleFunc("POST /admin/posts/new", requireAuth(handleAdminPostNewSubmit))
	mux.HandleFunc("GET /admin/posts/{id}/edit", requireAuth(handleAdminPostEditForm))
	mux.HandleFunc("POST /admin/posts/{id}/edit", requireAuth(handleAdminPostEditSubmit))
	mux.HandleFunc("POST /admin/posts/{id}/delete", requireAuth(handleAdminPostDelete))

	log.Println("Server starting on :5000")
	log.Fatal(http.ListenAndServe(":5000", mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleProjects(w http.ResponseWriter, r *http.Request) {
	page := 1
	if raw := r.URL.Query().Get("page"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			http.Error(w, "page must be a positive integer", http.StatusBadRequest)
			return
		}
		page = n
	}

	result, err := getGithubProjects(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getGithubProjects(page int) (ProjectsPage, error) {
	result := ProjectsPage{Projects: []Project{}, Page: page}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return result, nil
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: projectsPerPage, Page: page},
		Visibility:  "public",
		Affiliation: "owner",
		// "pushed" (last commit) rather than "updated" (any metadata change),
		// so the first page holds the most recently worked-on repositories.
		Sort:      "pushed",
		Direction: "desc",
	}

	repos, resp, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		return ProjectsPage{}, err
	}

	for _, repo := range repos {
		if !repo.GetPrivate() {
			result.Projects = append(result.Projects, getProjectData(repo))
		}
	}
	result.HasMore = resp != nil && resp.NextPage != 0

	return result, nil
}
