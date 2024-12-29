package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"linklite/config"
	"linklite/storage"
	"linklite/templates"
)

var store *storage.URLStore

func generateShortCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:6], nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		dynamicHandler(w, r)
		return
	}

	data := templates.IndexData{}

	switch r.Method {
	case http.MethodGet:
		templates.Index(data).Render(r.Context(), w)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			data.Error = "Invalid form data"
			templates.Index(data).Render(r.Context(), w)
			return
		}

		longURL := r.FormValue("url")
		if longURL == "" {
			data.Error = "URL is required"
			templates.Index(data).Render(r.Context(), w)
			return
		}

		// Validate URL
		if _, err := url.ParseRequestURI(longURL); err != nil {
			data.Error = "Invalid URL format"
			templates.Index(data).Render(r.Context(), w)
			return
		}

		// Check if URL already exists
		if shortCode, exists := store.FindByURL(longURL); exists {
			data.ShortenedURL = fmt.Sprintf("http://%s/%s", r.Host, shortCode)
			templates.Index(data).Render(r.Context(), w)
			return
		}

		// Generate short code
		shortCode, err := generateShortCode()
		if err != nil {
			data.Error = "Error generating short URL"
			templates.Index(data).Render(r.Context(), w)
			return
		}

		// Save to store
		if err := store.Set(shortCode, longURL); err != nil {
			data.Error = "Error saving URL"
			templates.Index(data).Render(r.Context(), w)
			return
		}

		// Create the full shortened URL
		data.ShortenedURL = fmt.Sprintf("http://%s/%s", r.Host, shortCode)
		templates.Index(data).Render(r.Context(), w)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the path without leading and trailing slashes
	shortCode := strings.Trim(r.URL.Path, "/")

	// Look up the URL
	if longURL, exists := store.Get(shortCode); exists {
		http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		return
	}

	// URL not found - render 404 page
	w.WriteHeader(http.StatusNotFound)
	templates.NotFound(shortCode).Render(r.Context(), w)
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize store
	store, err = storage.NewURLStore(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}
	defer store.Close()

	http.HandleFunc("/", homeHandler)

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("LinkLite URL Shortener starting on http://%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
