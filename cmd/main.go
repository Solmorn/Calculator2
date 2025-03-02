package main

import (
	agent "internal/agent"
	orch "internal/orch"
	"time"
)

func main() {
	go func() {
		orch.Run()
	}()

	time.Sleep(1 * time.Second)

	go func() {
		agent.Run()
	}()
}
