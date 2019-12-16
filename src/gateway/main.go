package main

import (
	"log"
	"os"
	"os/exec"
)

func buildCmd(name string, cmd ...string) *exec.Cmd {
	c := exec.Command(name, cmd...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func main() {
	// Generates service handlers dynamically
	gen := buildCmd("go", "run", "hack/genproto/main.go")
	if err := gen.Run(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
	}
	format := buildCmd("go", "fmt", "internal/server/protos.go")
	if err := format.Run(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
	}
	server := buildCmd("go", "run", "cmd/gateway/main.go")
	if err := server.Run(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
	}
}
