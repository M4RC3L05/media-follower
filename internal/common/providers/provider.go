package providers

import (
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
)

type ProviderName string

const ITUNES_MUSIC_RELEASES_PROVIDER ProviderName = "ITUNES_MUSIC_RELEASES_PROVIDER"

var PROVIDERS []ProviderName = []ProviderName{
	ITUNES_MUSIC_RELEASES_PROVIDER,
}

type IReleaseProvider interface {
	LookupReleaseInput(term string) (*model.Inputs, error)
	FetchReleases(input model.Inputs) ([]model.Releases, error)
	Name() ProviderName
}
