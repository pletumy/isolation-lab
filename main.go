package main

import (
	"fmt"
	"log"
	"os"

	"isolation-lab/scenario"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run main.go [locks|read_committed|repeatable_read|serializable]")
	}

	switch os.Args[1] {
	case "locks":
		scenario.RunLocks()
	default:
		fmt.Println("Unknown scenario")
	}
}
