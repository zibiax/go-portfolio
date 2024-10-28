package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImageURL    string `json:"image_url"`
}

func main() {
	// Log the GITHUB_TOKEN (first few characters)
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("Warning: GITHUB_TOKEN is not set")
	} else {
		log.Printf("GITHUB_TOKEN is set (starts with: %s...)", token[:5])
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/projects", handleProjects)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling index request")
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling projects request")
	
	projects, err := getGithubProjects()
	if err != nil {
		log.Printf("Error getting projects: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d projects", len(projects))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		log.Printf("Error encoding projects: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getGithubProjects() ([]Project, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Add debugging for API call
	log.Println("Fetching repositories from GitHub...")
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Type:        "owner",
	}
	
	repos, _, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		log.Printf("GitHub API error: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d repositories from GitHub", len(repos))

	var projects []Project
	for _, repo := range repos {
		if repo.GetPrivate() {
			continue // Skip private repositories
		}
		projects = append(projects, Project{
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			URL:         repo.GetHTMLURL(),
			ImageURL:    repo.GetOwner().GetAvatarURL(),
		})
	}

	return projects, nil
}
