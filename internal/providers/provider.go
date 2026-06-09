package providers

import "github.com/m4rc3l05/media-follower/.gen/go-jet/model"

type EProviders string

type IReleaseProvider[I any, O any] interface {
	FetchReleases(input I) ([]O, error)

	Name() string

	FromReleaseToPersistance(inputPersistance model.Inputs, release O) (*model.Releases, error)
	FromPersistanceToInput(inputPersistance model.Inputs) (*I, error)
}
