package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"portwatch/internal/daemon"
)

const defaultConfig = "portwatch.json"

func main() {
	configPath := flag.String("config", defaultConfig, "path to portwatch config file")
	flag.Parse()

	d, err := daemon.New(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config %q: %v\n", *configPath, err)
		os.Exit(1)
	}

	if err := d.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "daemon error: %v\n", err)
		os.Exit(1)
	}
}
