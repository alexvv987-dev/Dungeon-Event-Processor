package main

import (
	"fmt"
	"os"

	"impulse/internal/config"
	"impulse/internal/engine"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: impulse <config.json> <events>")
		os.Exit(1)
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	eventsFile, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot open events:", err)
		os.Exit(1)
	}
	defer eventsFile.Close()

	out, err := engine.Process(cfg, eventsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(out)
}
