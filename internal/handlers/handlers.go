package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Solmorn/Calculator2/internal/patterns"
	calculator "github.com/Solmorn/Calculator2/pkg"
)

var (
	expressions = make(map[string]patterns.Expression)
	tasks       = make(map[string]patterns.Task)
	mu          sync.Mutex
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	id := generateID()
	mu.Lock()
	expressions[id] = patterns.Expression{ID: id, Status: "pending", Result: 0}
	mu.Unlock()

	go func() {
		time.Sleep(1 * time.Second)
		result, err := calculator.Calc(req.Expression)
		mu.Lock()
		if err != nil {
			expressions[id] = patterns.Expression{ID: id, Status: "error", Result: 0}
		} else {
			expressions[id] = patterns.Expression{ID: id, Status: "done", Result: result}
		}
		mu.Unlock()
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()

	expressionsList := make([]patterns.Expression, 0, len(expressions))
	for _, expr := range expressions {
		expressionsList = append(expressionsList, expr)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": expressionsList,
	})
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Path[len("/api/v1/expressions/"):]
	mu.Lock()
	defer mu.Unlock()
	if expr, ok := expressions[id]; ok {
		json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "expression not found"})
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()
	for _, task := range tasks {
		if task.Operation != "" {
			delete(tasks, task.ID)
			json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "no tasks available"})
}

func TaskResultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result patterns.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if expr, ok := expressions[result.ID]; ok {
		expr.Result = result.Result
		expr.Status = "done"
		expressions[result.ID] = expr
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "result updated"})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "expression not found"})
	}
}

func generateID() string {
	return "id" + time.Now().Format("20060102150405")
}
