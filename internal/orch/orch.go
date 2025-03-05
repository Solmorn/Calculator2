package orch

import (
	"fmt"
	"net/http"

	"github.com/Solmorn/Calculator2/internal/handlers"
)

func Run() {
	http.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)
	http.HandleFunc("/api/v1/expressions", handlers.ExpressionsHandler)

	http.HandleFunc("/api/v1/expressions/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/api/v1/expressions/{id}"
		r.URL.RawPath = "/api/v1/expressions/{id}"
		handlers.ExpressionHandler(w, r)
	})

	http.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handlers.TaskHandler(w, r)
		} else if r.Method == "POST" {
			handlers.TaskResultHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Orchestrator running on :8080")
	http.ListenAndServe(":8080", nil)
}
