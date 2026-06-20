package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
)

type ItunesMusicReleaseProvider struct {
	Validator *validator.Validate
	Conform   *mold.Transformer
}

var (
	_            IReleaseProvider = ItunesMusicReleaseProvider{}
	ogImageRegex                  = regexp.MustCompile(
		`(?i)<meta\s+property="og:image"\s+content="([^"]*)"`,
	)
)

type ItunesLookupArtist struct {
	AmgArtistId      *int64 `json:"amgArtistId,omitempty" validate:"omitempty"`
	ArtistId         int64  `json:"artistId"              validate:"required"`
	ArtistLinkUrl    string `json:"artistLinkUrl"         validate:"required,url"`
	ArtistName       string `json:"artistName"            validate:"required"`
	ArtistType       string `json:"artistType"            validate:"required,eq=Artist"`
	PrimaryGenreId   int64  `json:"primaryGenreId"        validate:"required"`
	PrimaryGenreName string `json:"primaryGenreName"      validate:"required"`
	WrapperType      string `json:"wrapperType"           validate:"required,eq=artist"`
}

type ItunesLookupAlbum struct {
	AmgArtistID            *int64   `json:"amgArtistId,omitempty"           validate:"omitempty"`
	ArtistID               int64    `json:"artistId"                        validate:"required"`
	ArtistName             string   `json:"artistName"                      validate:"required"`
	ArtistViewURL          *string  `json:"artistViewUrl,omitempty"         validate:"omitempty,url"`
	ArtworkURL100          string   `json:"artworkUrl100"                   validate:"required,url"`
	ArtworkURL60           string   `json:"artworkUrl60"                    validate:"required,url"`
	CollectionCensoredName string   `json:"collectionCensoredName"          validate:"required"`
	CollectionExplicitness string   `json:"collectionExplicitness"          validate:"required"`
	CollectionID           int64    `json:"collectionId"                    validate:"required"`
	CollectionName         string   `json:"collectionName"                  validate:"required"`
	CollectionPrice        *float64 `json:"collectionPrice,omitempty"       validate:"omitempty"`
	CollectionType         string   `json:"collectionType"                  validate:"required,eq=Album"`
	CollectionViewURL      string   `json:"collectionViewUrl"               validate:"required,url"`
	ContentAdvisoryRating  *string  `json:"contentAdvisoryRating,omitempty" validate:"omitempty"`
	Copyright              *string  `json:"copyright,omitempty"             validate:"omitempty"`
	Country                string   `json:"country"                         validate:"required"`
	Currency               string   `json:"currency"                        validate:"required"`
	PrimaryGenreName       string   `json:"primaryGenreName"                validate:"required"`
	ReleaseDate            *string  `json:"releaseDate,omitempty"           validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" mod:"toisoutc"`
	TrackCount             int64    `json:"trackCount"                      validate:"required"`
	WrapperType            string   `json:"wrapperType"                     validate:"required,eq=collection"`
}

type ItunesResponse[T any] struct {
	ResultCount int `json:"resultCount" validate:"gte=0"`
	Results     []T `json:"results"     validate:"required,dive" mod:"dive"`
}

func (i ItunesMusicReleaseProvider) fetchAlbums(
	id string,
) (_ *ItunesResponse[ItunesLookupAlbum], eOut error) {
	response, err := http.Get(
		fmt.Sprintf(
			"https://itunes.apple.com/lookup?id=%s&entity=album&media=music&sort=recent&limit=60",
			id,
		),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			if eOut != nil {
				eOut = errors.Join(eOut, err)
			} else {
				eOut = err
			}
		}
	}()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data ItunesResponse[ItunesLookupAlbum]
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Results) > 0 {
		data.Results = slices.Delete(data.Results, 0, 1)
	}

	err = i.Conform.Struct(context.Background(), &data)
	if err != nil {
		return nil, err
	}

	err = i.Validator.Struct(data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func resolveArtistImage(url string) (u *string, outError error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			if outError != nil {
				outError = errors.Join(outError, err)
			} else {
				outError = err
			}

			u = nil
		}
	}()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	html := string(data)

	matches := ogImageRegex.FindStringSubmatch(html)
	if len(matches) < 2 {
		return nil, nil
	}

	imageURL := matches[1]
	if imageURL == "" || strings.Contains(imageURL, "apple-music-") {
		return nil, nil
	}

	parts := strings.Split(imageURL, "/")
	lastPart := parts[len(parts)-1]

	dotIdx := strings.LastIndex(lastPart, ".")
	if dotIdx == -1 {
		return nil, nil
	}

	ext := lastPart[dotIdx+1:]
	parts[len(parts)-1] = fmt.Sprintf("300x300.%s", ext)
	final := strings.Join(parts, "/")
	return &final, nil
}

func (i ItunesMusicReleaseProvider) LookupReleaseInput(
	term string,
) (input *model.Inputs, outError error) {
	response, err := http.Get(fmt.Sprintf("https://itunes.apple.com/lookup?id=%s", term))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			if outError != nil {
				outError = errors.Join(outError, err)
			} else {
				outError = err
			}

			input = nil
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var parsed ItunesResponse[ItunesLookupArtist]
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	if err := i.Validator.Struct(parsed); err != nil {
		return nil, err
	}

	if len(parsed.Results) <= 0 {
		return nil, fmt.Errorf("search term \"%s\" did not return any input", term)
	}
	artist := parsed.Results[0]
	artistImage, err := resolveArtistImage(artist.ArtistLinkUrl)
	if err != nil {
		return nil, err
	}

	return &model.Inputs{
		InternalProviderID: strconv.FormatInt(artist.ArtistId, 10),
		Provider:           string(i.Name()),
		Name:               artist.ArtistName,
		ImageURL:           artistImage,
		ExternalLink:       &artist.ArtistLinkUrl,
	}, nil
}

func imgUrl(in string) string {
	splited := strings.Split(in, "/")
	splited[len(splited)-1] = "300x300bb.jpg"

	return strings.Join(splited, "/")
}

func (i ItunesMusicReleaseProvider) FetchReleases(input model.Inputs) ([]model.Releases, error) {
	albums, err := i.fetchAlbums(input.InternalProviderID)
	if err != nil {
		return nil, err
	}

	albums.Results = slices.DeleteFunc(albums.Results, func(album ItunesLookupAlbum) bool {
		isCompilation := strings.Contains(strings.ToLower(album.ArtistName), "varios artists")
		isDJMix := strings.Contains(strings.ToLower(album.CollectionName), "dj mix") ||
			strings.Contains(strings.ToLower(album.CollectionCensoredName), "dj mix")
		noReleaseInfo := album.ReleaseDate == nil

		return isCompilation || isDJMix || noReleaseInfo
	})

	final := make([]model.Releases, len(albums.Results))
	for k, album := range albums.Results {
		img := imgUrl(album.ArtworkURL100)
		final[k] = model.Releases{
			InternalProviderID: strconv.FormatInt(album.CollectionID, 10),
			InputID:            input.ID,
			Title:              fmt.Sprintf("%s - %s", album.ArtistName, album.CollectionName),
			ImageURL:           &img,
			ReleasedAt:         *album.ReleaseDate,
			ExternalLink:       &album.CollectionViewURL,
		}
	}

	return final, nil
}

func (i ItunesMusicReleaseProvider) Name() ProviderName {
	return ITUNES_MUSIC_RELEASES_PROVIDER
}
