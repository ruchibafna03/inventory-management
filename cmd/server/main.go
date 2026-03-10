package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/val/inventory/internal/api"
	"github.com/val/inventory/internal/db"
)

func main() {
	_ = godotenv.Load()

	database := db.Connect()
	defer database.Close()

	runMigrations(database)

	router := api.NewRouter(database)

	// Serve React SPA from WEB_ROOT (set in production, defaults to ./web)
	webRoot := os.Getenv("WEB_ROOT")
	if webRoot == "" {
		webRoot = filepath.Join(".", "web")
	}

	// Mount SPA handler for everything that isn't /api
	router.Get("/*", spaHandler(webRoot))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("VAL Inventory running on http://0.0.0.0:%s  (web=%s)\n", port, webRoot)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// spaHandler serves static files from webRoot and falls back to index.html
// for any path that doesn't match a real file (React Router support).
func spaHandler(webRoot string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(webRoot))
	return func(w http.ResponseWriter, r *http.Request) {
		// Strip the chi wildcard prefix
		path := chi.URLParam(r, "*")
		if path == "" {
			path = "index.html"
		}

		// If request is for a known asset extension, serve directly
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".js" || ext == ".css" || ext == ".png" || ext == ".svg" ||
			ext == ".ico" || ext == ".woff" || ext == ".woff2" || ext == ".ttf" {
			fs.ServeHTTP(w, r)
			return
		}

		// Check if the file exists on disk
		fullPath := filepath.Join(webRoot, filepath.Clean("/"+path))
		if _, err := os.Stat(fullPath); err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routes
		http.ServeFile(w, r, filepath.Join(webRoot, "index.html"))
	}
}
