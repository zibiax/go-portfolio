package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
    "context"

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
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for _, repo := range repos {
		projects = append(projects, Project{
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			URL:         repo.GetHTMLURL(),
			ImageURL:    repo.GetOwner().GetAvatarURL(),
		})
	}

	return projects, nil
}
