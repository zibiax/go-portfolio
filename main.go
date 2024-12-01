package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v35/github"
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
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("Warning: GITHUB_TOKEN is not set")
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/projects", handleProjects)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
	projects, err := getGithubProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func getGithubProjects() ([]Project, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, nil
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Type:        "owner",
		Sort:        "updated",
		Direction:   "desc",
	}
	
	repos, _, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for _, repo := range repos {
		if !repo.GetPrivate() {
			projects = append(projects, getProjectData(repo))
		}
	}

	return projects, nil
}
