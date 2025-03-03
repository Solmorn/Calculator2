package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Solmorn/Calculator2/internal/patterns"
)

func Run() {

	computingPowerStr := os.Getenv("COMPUTING_POWER")
	computingPower, err := strconv.Atoi(computingPowerStr)
	if err != nil || computingPower <= 0 {
		computingPower = 4
	}

	var wg sync.WaitGroup
	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, err := fetchTask()
		if err != nil {
			log.Printf("Worker %d: failed to fetch task: %v\n", id, err)
			time.Sleep(1 * time.Second)
			continue
		}

		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		result := compute(task)

		err = postResult(task.ID, result)
		if err != nil {
			log.Printf("Worker %d: failed to send result: %v\n", id, err)
		} else {
			log.Printf("Worker %d: successfully processed task %s\n", id, task.ID)
		}
	}
}

func fetchTask() (*patterns.Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no available tasks or unexpected status code: %d", resp.StatusCode)
	}

	var taskResponse struct {
		Task patterns.Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&taskResponse); err != nil {
		return nil, fmt.Errorf("failed to decode task: %v", err)
	}

	return &taskResponse.Task, nil
}

func compute(task *patterns.Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			log.Printf("Division by zero!")
			return 0
		}
		return task.Arg1 / task.Arg2
	default:
		log.Printf("Unknown operation: %s", task.Operation)
		return 0
	}
}

func postResult(taskID string, result float64) error {
	taskResult := patterns.TaskResult{
		ID:     taskID,
		Result: result,
	}

	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		return fmt.Errorf("failed to marshal task result: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send result: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
