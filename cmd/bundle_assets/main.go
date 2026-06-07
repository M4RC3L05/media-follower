package main

import (
	"flag"
	"os"
	"strings"

	"github.com/m4rc3l05/media-follower/internal/commands"
)

type Entrypoints []string

func (i Entrypoints) String() string {
	return strings.Join(i, ":")
}

func (i *Entrypoints) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func run(entries []string, outDir string) int {
	if errs := commands.BundleAssets(entries, outDir); errs != nil {
		return 1
	}

	return 0
}

func main() {
	var entries Entrypoints
	flag.Var(
		&entries,
		"e",
		"Entry directory globs where your assets live, can be specified multiple times",
	)
	outDir := flag.String("o", "", "Output directory where you bundle assets will live")

	flag.Parse()

	os.Exit(run(entries, *outDir))
}
