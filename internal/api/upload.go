package api

import (
	"encoding/json"
	"gutondeka/internal/db"
	"gutondeka/internal/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Only POST method is allowed"}`))
		return
	}

	// Parse multipart form data (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse form data"})
		return
	}

	// Get article metadata
	contributorIDStr := r.FormValue("contributor_id")
	contributorID, err := strconv.Atoi(contributorIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid contributor_id: must be a number"})
		return
	}

	// Check if contributor exists
	exists, err := db.CheckContributorExists(contributorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to verify contributor"})
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Contributor not found"})
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	source := r.FormValue("source")
	category := r.FormValue("category")
	dateStr := r.FormValue("date")

	if title == "" || author == "" || source == "" || category == "" || dateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing required fields"})
		return
	}

	// Parse date
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format. Use RFC3339 format (e.g., 2024-03-20T12:00:00Z)"})
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	exists, err = db.CheckFileExists(handler.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to check file existence"})
		return
	}
	if exists {
		w.WriteHeader(http.StatusAlreadyReported)
		json.NewEncoder(w).Encode(map[string]string{"error": "File already exists"})
		return
	}

	// Use original filename
	filePath := filepath.Join("data", "articles", handler.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
		return
	}

	article := models.Article{
		ContributorID: contributorID,
		Title:         title,
		Author:        author,
		Source:        source,
		Date:          date,
		Category:      category,
		FilePath:      filePath,
	}

	if err := db.InsertArticle(&article); err != nil {
		os.Remove(filePath)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save article"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Article uploaded successfully",
		"file_path": filePath,
	})
}
