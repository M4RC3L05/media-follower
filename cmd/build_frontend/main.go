package main

import (
	"os"

	"github.com/m4rc3l05/media-follower/internal/commands"
)

func run(entry []string, dist string) int {
	if err := commands.BuildFrontend(entry, dist); err != nil {
		return -1
	}

	return 0
}

func main() {
	entries := os.Args[1:max(len(os.Args)-1, 0)]
	dist := os.Args[max(len(os.Args)-1, 0)]

	os.Exit(run(entries, dist))
}
