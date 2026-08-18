package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"html/template"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const sessionCookieName = "session"
const sessionDuration = 7 * 24 * time.Hour

func generateSessionToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func createSession() (string, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	_, err = db.Exec(
		`INSERT INTO sessions (token_hash, created_at, expires_at) VALUES (?, ?, ?)`,
		hashToken(token), now, now.Add(sessionDuration),
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func isValidSession(token string) bool {
	if token == "" {
		return false
	}
	var expiresAt time.Time
	err := db.QueryRow(`SELECT expires_at FROM sessions WHERE token_hash = ?`, hashToken(token)).Scan(&expiresAt)
	if err != nil {
		return false
	}
	return time.Now().UTC().Before(expiresAt)
}

func deleteSession(token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE token_hash = ?`, hashToken(token))
	return err
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || !isValidSession(cookie.Value) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	renderLogin(w, "")
}

func renderLogin(w http.ResponseWriter, errMsg string) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{"Error": errMsg})
}

func handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	wantUsername := os.Getenv("ADMIN_USERNAME")
	wantPasswordHash := os.Getenv("ADMIN_PASSWORD_HASH")

	usernameMatches := subtle.ConstantTimeCompare([]byte(username), []byte(wantUsername)) == 1
	passwordMatches := wantPasswordHash != "" && bcrypt.CompareHashAndPassword([]byte(wantPasswordHash), []byte(password)) == nil

	if !usernameMatches || !passwordMatches {
		renderLogin(w, "Invalid username or password")
		return
	}

	token, err := createSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(sessionDuration),
	})

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		deleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
