package main

import (
	json2 "encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"jira-parser/jira_fetch"
	"net/http"
	"strings"
)

func setUpRoute() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value for preflight requests cache in seconds
	}))
	router.HandleFunc("/jira/{jira_id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "jira_id")
		resp, err := jira_fetch.Fetch(id)
		cleanedText := strings.Replace(resp, "```json", "", -1)
		cleanedText = strings.Replace(cleanedText, "```", "", -1)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json2.NewEncoder(w).Encode(cleanedText)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("response sent.... ")

	})
	return router
}

func main() {
	router := setUpRoute()

	fmt.Println("Starting server...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}

}
