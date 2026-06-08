package main

import (
	"flag"
	"os"
	"strings"

	"github.com/m4rc3l05/media-follower/internal/commands"
)

func run(entries []string, outDir string) int {
	if errs := commands.BundleAssets(entries, outDir); errs != nil {
		return 1
	}

	return 0
}

func main() {
	entries := flag.String(
		"e",
		"",
		"Entry directory globs where your assets live, coma seperated list of globs",
	)
	outDir := flag.String("o", "", "Output directory where you bundle assets will live")

	flag.Parse()

	os.Exit(run(strings.Split(*entries, ","), *outDir))
}
