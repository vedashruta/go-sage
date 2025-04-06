package router

import (
	"net/http"
	"server/modules/search"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	search.RegisterRoutes(mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	return mux
}
