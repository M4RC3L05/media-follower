package commands

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

var httpImportsPlugin = api.Plugin{
	Name: "http-imports",
	Setup: func(build api.PluginBuild) {
		httpImportsCacheDir := ".cache/http-imports"

		fetchHttpImport := func(p string) (_ []byte, errOut error) {
			hash := md5.New()
			if _, err := io.WriteString(hash, p); err != nil {
				return nil, err
			}
			md5Str := hex.EncodeToString(hash.Sum(nil))
			finalPath := path.Join(httpImportsCacheDir, md5Str)

			_, err := os.Stat(finalPath)
			if err == nil {
				return os.ReadFile(finalPath)
			}

			res, err := http.Get(p)
			if err != nil {
				return nil, err
			}

			defer func() {
				if err := res.Body.Close(); err != nil {
					errOut = err
				}
			}()

			contents, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			if err := os.WriteFile(finalPath, contents, 0o644); err != nil {
				return nil, err
			}

			return os.ReadFile(finalPath)
		}

		resolveLoaderFrompPath := func(path string) (api.Loader, error) {
			if strings.HasSuffix(path, ".css") {
				return api.LoaderCSS, nil
			}

			if strings.HasSuffix(path, ".js") {
				return api.LoaderJS, nil
			}

			return api.LoaderEmpty, fmt.Errorf("could not resolve loader for \"%s\"", path)
		}

		build.OnStart(func() (api.OnStartResult, error) {
			if err := os.MkdirAll(httpImportsCacheDir, 0o755); err != nil {
				return api.OnStartResult{}, err
			}

			return api.OnStartResult{}, nil
		})

		build.OnResolve(api.OnResolveOptions{Filter: `^https://`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      args.Path,
					Namespace: "http-imports",
				}, nil
			})

		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "http-imports"},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				loader, err := resolveLoaderFrompPath(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				res, err := fetchHttpImport(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				contents := string(res)

				return api.OnLoadResult{
					Contents: &contents,
					Loader:   loader,
				}, nil
			})
	},
}

func BundleAssets(entries []string, outDir string) *[]api.Message {
	result := api.Build(api.BuildOptions{
		LogLevel:          api.LogLevelInfo,
		EntryPoints:       entries,
		Bundle:            true,
		Outdir:            outDir,
		Plugins:           []api.Plugin{httpImportsPlugin},
		Write:             true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Target:            api.ESNext,
		Format:            api.FormatESModule,
		Platform:          api.PlatformBrowser,
		TreeShaking:       api.TreeShakingTrue,
		Splitting:         true,
	})

	if len(result.Errors) > 0 {
		return &result.Errors
	}

	return nil
}
