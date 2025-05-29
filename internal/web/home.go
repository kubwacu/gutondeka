package web

import (
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/web/templates/home.html")
}
