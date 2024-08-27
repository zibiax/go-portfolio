package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()

    // Serve the main page
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    // HTMX endpoints
    r.GET("/about", func(c *gin.Context) {
        c.HTML(http.StatusOK, "about.html", nil)
    })
    r.GET("/projects", func(c *gin.Context) {
        c.HTML(http.StatusOK, "projects.html", nil)
    })
    r.GET("/contact", func(c *gin.Context) {
        c.HTML(http.StatusOK, "contact.html", nil)
    })

    // Serve robots.txt
    r.StaticFile("/robots.txt", "./robots.txt")

    // Static assets like CSS and JavaScript
    r.Static("/assets", "./assets")

    // Run the server
    r.Run(":8080")
}

