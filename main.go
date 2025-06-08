package main

import (
	"fmt"
	"gutondeka/internal/api"
	"gutondeka/internal/db"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

func handleRequests() {
	http.HandleFunc("/api/upload", enableCORS(api.UploadHandler))
	http.HandleFunc("/api/overview", enableCORS(api.OverviewHandler))
}

func main() {
	articlesDir := filepath.Join("data", "articles")
	if err := os.MkdirAll(articlesDir, 0755); err != nil {
		log.Fatal("Failed to create articles directory:", err)
	}

	if err := db.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	handleRequests()

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
