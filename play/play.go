package main

import (
	"fmt"

	"github.com/PaulioRandall/go-trackerr"
	"github.com/PaulioRandall/go-trackerr/play/example"
)

func init() {
	trackerr.Initialised()
}

func main() {
	fmt.Println("Play...")
	example.Run()
}
