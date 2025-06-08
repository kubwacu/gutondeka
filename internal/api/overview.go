package api

import (
	"encoding/json"
	"net/http"

	"gutondeka/internal/db"
	"gutondeka/internal/models"
)

func OverviewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Only GET method is allowed"}`))
		return
	}

	articleCounts, err := db.GetArticlesCount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to get article counts"}`))
		return
	}

	defaultEntry := models.ArticleCount{
		MonthYear: "05-2025",
		Count:     0,
	}

	finalCounts := []models.ArticleCount{defaultEntry}
	for _, count := range articleCounts {
		if count.MonthYear != "05-2025" {
			finalCounts = append(finalCounts, models.ArticleCount{
				MonthYear: count.MonthYear,
				Count:     count.Count,
			})
		}
	}

	response := map[string]interface{}{
		"article_counts": finalCounts,
		"total_articles": 0,
	}

	for _, count := range finalCounts {
		response["total_articles"] = response["total_articles"].(int) + count.Count
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to generate response"}`))
		return
	}

	w.Write(jsonResponse)
}
