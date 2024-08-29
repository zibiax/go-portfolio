package main

import (
    "net/http"
    "bytes"

    "github.com/a-h/templ"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    templateEngine, err := templ.New()
    if err != nil {
        panic(err) // Handle initialization error
    }

    // Serve the main page
    r.GET("/", func(c *gin.Context) {
        var buf bytes.Buffer
        err := templateEngine.RenderTemplate(&buf, "templates/index.templ", gin.H{
            "title": "Home Page",
        })
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }
        c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
    })

    // HTMX endpoints
    setupHTMXEndpoints(r, templateEngine)

    // Serve robots.txt and static assets
    r.StaticFile("/robots.txt", "./robots.txt")
    r.Static("/assets", "./assets")

    // Run the server
    r.Run(":8080")
}

func setupHTMXEndpoints(r *gin.Engine, templateEngine *templ.Template) {
    // endpoints to template file
    endpoints := map[string]string{
        "/about": "templates/about.templ",
        "/projects": "templates/projects.templ",
        "/contact": "templates/contact.templ",
    }

    for path, tmpl := range endpoints {
        r.GET(path, func(c *gin.Context) {
            var buf bytes.Buffer
            err := templateEngine.RenderTemplate(&buf, tmpl, gin.H{
                "title": c.Param("title"),
            })
            if err != nil {
                c.AbortWithError(http.StatusInternalServerError, err)
                return
            }
            c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
        })
    }
}


