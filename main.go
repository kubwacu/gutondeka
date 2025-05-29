package main

import (
	"fmt"
	"gutondeka/internal/api"
	"gutondeka/internal/web"
	"log"
	"net/http"
)

func handleRequests() {
	http.HandleFunc("/", web.HomePageHandler)
	http.HandleFunc("/api/upload", api.UploadHandler)
	http.HandleFunc("/api/exists", api.ExistsHandler)
}

func main() {
	handleRequests()

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
