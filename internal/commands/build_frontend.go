package commands

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func resolveLoader(path string) api.Loader {
	if strings.HasSuffix(path, ".css") {
		return api.LoaderCSS
	}

	return api.LoaderNone
}

var httpImportsPlugin = api.Plugin{
	Name: "http_imports",
	Setup: func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `^https?://`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      args.Path,
					Namespace: "http_imports_ns",
				}, nil
			})

		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "http_imports_ns"},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				response, err := http.Get(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				defer response.Body.Close() //nolint:errcheck

				contents, err := io.ReadAll(response.Body)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				return api.OnLoadResult{
					Contents: new(string(contents)),
					Loader:   resolveLoader(args.Path),
				}, nil
			})
	},
}

func BuildFrontend(entry []string, dest string) error {
	if err := os.RemoveAll("./internal/entrypoints/admin_server/dist"); err != nil {
		return err
	}

	result := api.Build(api.BuildOptions{
		LogLevel:          api.LogLevelInfo,
		EntryPoints:       entry,
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Splitting:         true,
		Target:            api.ESNext,
		TreeShaking:       api.TreeShakingTrue,
		Platform:          api.PlatformBrowser,
		Format:            api.FormatESModule,
		Loader: map[string]api.Loader{
			".ico":     api.LoaderCopy,
			".gitkeep": api.LoaderCopy,
		},
		Outdir:  dest,
		Plugins: []api.Plugin{httpImportsPlugin},
		Write:   true,
	})

	if len(result.Errors) != 0 {
		return errors.New("could not bundle")
	}

	return nil
}
