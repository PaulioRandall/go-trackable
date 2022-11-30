package main

import (
	"log"
	"os"

	"github.com/PaulioRandall/trackable/play/examples"
)

func main() {
	if len(os.Args) < 2 {
		runExample("1")
	} else {
		runExample(os.Args[1])
	}
}

func runExample(id string) {
	switch id {
	case "1", "workflow":
		examples.Workflow()
	case "2", "packages":
		examples.Packages()
	default:
		log.Fatalf("Unknown example ID %q", id)
	}
}
