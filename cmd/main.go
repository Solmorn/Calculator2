package main

import (
	"time"

	agent "github.com/Solmorn/Calculator2/internal/agent"
	orch "github.com/Solmorn/Calculator2/internal/orch"
)

func main() {
	go func() {
		orch.Run()
	}()

	time.Sleep(1 * time.Second)

	go func() {
		agent.Run()
	}()

	select {}
}
