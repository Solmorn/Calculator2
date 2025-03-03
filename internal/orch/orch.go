package orch

import (
	"encoding/json"

	"log"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	Result        float64 `json:"result,omitempty"`
	Status        string  `json:"status"`
	OperationTime int     `json:"operation_time"`
}

var (
	tasks         = make(map[int]*Task)
	taskIDCounter = 1
	taskMutex     = &sync.Mutex{}
)

func Run() {
	http.HandleFunc("/api/v1/calculate", handleCalculate)
	http.HandleFunc("/api/v1/expressions", handleExpressions)
	http.HandleFunc("/api/v1/expressions/", handleExpressionByID)
	http.HandleFunc("/internal/task", handleGetTask)
	http.HandleFunc("/internal/result", handleResult)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	task := &Task{
		ID:            taskIDCounter,
		Status:        "pending",
		OperationTime: 1000, // Example time, should be read from env vars
	}
	taskIDCounter++

	taskMutex.Lock()
	tasks[task.ID] = task
	taskMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": task.ID})
}

func handleExpressions(w http.ResponseWriter, r *http.Request) {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	exprList := make([]map[string]interface{}, 0, len(tasks))
	for _, t := range tasks {
		exprList = append(exprList, map[string]interface{}{
			"id":     t.ID,
			"status": t.Status,
			"result": t.Result,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": exprList,
	})
}

func handleExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	taskID, _ := strconv.Atoi(id)

	taskMutex.Lock()
	task, exists := tasks[taskID]
	taskMutex.Unlock()

	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]*Task{"expression": task})
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	for _, task := range tasks {
		if task.Status == "pending" {
			task.Status = "in_progress"
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "No tasks available", http.StatusNotFound)
}

func handleResult(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	taskMutex.Lock()
	task, exists := tasks[result.ID]
	taskMutex.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	taskMutex.Lock()
	task.Result = result.Result
	task.Status = "completed"
	taskMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}
