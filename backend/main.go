package main

import (
	"fmt"
	"log"
	"net/http"
	"server/config"
	"server/driver/csv"
	searchengine "server/driver/search-engine"
	"server/router"
	"time"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	searchengine.Init()
	dirPath := "storage/upload"
	// err = parquet.Parse(dirPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err = csv.ParseCSVFromDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	handlerWithCORS := withCORS(router.NewRouter())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port),
		Handler:      handlerWithCORS,
		ReadTimeout:  time.Duration(config.Config.ReadTimeout),
		WriteTimeout: time.Duration(config.Config.WriteTimeout),
	}
	fmt.Printf("Server running on %s:%d\n", config.Config.Host, config.Config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins - for dev only
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
