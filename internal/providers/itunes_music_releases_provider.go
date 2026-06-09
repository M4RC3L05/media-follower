package providers

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/internal/providers/inputs"
	"github.com/m4rc3l05/media-follower/internal/providers/outputs"
)

type ItunesMusicReleasesProvider struct {
	InputProvider  inputs.IInputProvider[inputs.ItunesArtist]
	OutputProvider outputs.IOutputProvider[inputs.ItunesArtist, outputs.ItunesAlbum]
}

var _ IReleaseProvider[inputs.ItunesArtist, outputs.ItunesAlbum] = ItunesMusicReleasesProvider{}

func NewItunesMusicReleasesProvider(
	inProvider inputs.IInputProvider[inputs.ItunesArtist],
	outProvider outputs.IOutputProvider[inputs.ItunesArtist, outputs.ItunesAlbum],
) ItunesMusicReleasesProvider {
	return ItunesMusicReleasesProvider{
		InputProvider:  inProvider,
		OutputProvider: outProvider,
	}
}

func (i ItunesMusicReleasesProvider) FetchReleases(
	input inputs.ItunesArtist,
) ([]outputs.ItunesAlbum, error) {
	releases, err := i.OutputProvider.FetchOutputs(input)
	if err != nil {
		return nil, err
	}

	// Lets ignore albums that:
	// - Are a compilation (have `various artists` on the artist name)
	// - Are a DJ Mix (have `various artists` on tthe collection [album] name)
	// - Do not have release date info
	releases = slices.DeleteFunc(releases, func(output outputs.ItunesAlbum) bool {
		isCompilation := strings.Contains(
			strings.ToLower(output.ArtistName),
			"various artists",
		)
		isDjMix := strings.Contains(
			strings.ToLower(output.CollectionCensoredName),
			"various artists",
		) ||
			strings.Contains(strings.ToLower(output.CollectionName), "various artists")
		noReleaseDate := output.ReleaseDate == nil

		return isCompilation || isDjMix || noReleaseDate
	})

	return releases, nil
}

func (i ItunesMusicReleasesProvider) FromReleaseToPersistance(
	inputPersistance model.Inputs,
	release outputs.ItunesAlbum,
) (*model.Releases, error) {
	encodedRaw, err := i.OutputProvider.JSONEncode(release)
	if err != nil {
		return nil, err
	}

	releasedAt, err := time.Parse(time.RFC3339, *release.ReleaseDate)
	if err != nil {
		return nil, err
	}

	return &model.Releases{
		ID:            strconv.FormatInt(release.CollectionID, 10),
		InputID:       inputPersistance.ID,
		InputProvider: inputPersistance.Provider,
		Provider:      i.Name(),
		Raw:           encodedRaw,
		ReleasedAt:    releasedAt.Format("2006-01-02T15:04:05.000Z07:00"),
	}, nil
}

func (i ItunesMusicReleasesProvider) FromPersistanceToInput(
	inputPersistance model.Inputs,
) (*inputs.ItunesArtist, error) {
	var itunesArtist inputs.ItunesArtist
	err := json.Unmarshal(inputPersistance.Raw, &itunesArtist)
	if err != nil {
		return nil, err
	}

	err = i.InputProvider.Validate(itunesArtist)
	if err != nil {
		return nil, err
	}

	return &itunesArtist, nil
}

func (i ItunesMusicReleasesProvider) Name() string {
	return "ITUNES_MUSIC_RELEASES_PROVIDER"
}
