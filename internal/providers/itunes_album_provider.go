package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/.gen/jetdb/model"
	"github.com/m4rc3l05/media-follower/internal/common"
)

const ITUNES_ALBUM_PROVIDER_NAME = "ITUNES_ALBUM_PROVIDER"

type ItunesAlbum struct {
	AmgArtistID            *int64   `json:"amgArtistId,omitempty"           validate:"omitempty,number"`
	ArtistID               int64    `json:"artistId"                        validate:"required,number"`
	ArtistName             string   `json:"artistName"                      validate:"required"`
	ArtistViewURL          string   `json:"artistViewUrl"                   validate:"required,url"`
	ArtworkURL100          string   `json:"artworkUrl100"                   validate:"required,url"`
	ArtworkURL60           string   `json:"artworkUrl60"                    validate:"required,url"`
	CollectionCensoredName string   `json:"collectionCensoredName"          validate:"required"`
	CollectionExplicitness string   `json:"collectionExplicitness"          validate:"required"`
	CollectionID           int64    `json:"collectionId"                    validate:"required,number"`
	CollectionName         string   `json:"collectionName"                  validate:"required"`
	CollectionPrice        *float64 `json:"collectionPrice,omitempty"       validate:"omitempty,number"`
	CollectionType         string   `json:"collectionType"                  validate:"required,eq=Album"`
	CollectionViewURL      string   `json:"collectionViewUrl"               validate:"required,url"`
	ContentAdvisoryRating  *string  `json:"contentAdvisoryRating,omitempty"`
	Copyright              *string  `json:"copyright,omitempty"`
	Country                string   `json:"country"                         validate:"required"`
	Currency               string   `json:"currency"                        validate:"required"`
	PrimaryGenreName       string   `json:"primaryGenreName"                validate:"required"`
	ReleaseDate            *string  `json:"releaseDate,omitempty"           validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TrackCount             int64    `json:"trackCount"                      validate:"required,number"`
	WrapperType            string   `json:"wrapperType"                     validate:"required,eq=collection"`
}

type ItunesAlbumProvider struct {
	validate *validator.Validate
	log      *slog.Logger
}

// Compile time check that providers implement interface
var (
	_ IOutputProvider[ItunesArtist, ItunesAlbum] = ItunesAlbumProvider{}
)

func NewItunesAlbumProvider(validator *validator.Validate) ItunesAlbumProvider {
	return ItunesAlbumProvider{
		validate: validator,
		log:      common.NewLogger("itunes-album-releases-provider"),
	}
}

func fetchAlbums(id int64) (_ *ItunesResponseModel[ItunesAlbum], eOut error) {
	response, err := http.Get(
		fmt.Sprintf(
			"https://itunes.apple.com/lookup?id=%d&entity=album&media=music&sort=recent&limit=60",
			id,
		),
	)

	defer func() {
		if response != nil && response.Body != nil {
			if err := response.Body.Close(); err != nil {
				if eOut != nil {
					eOut = errors.Join(eOut, err)
				} else {
					eOut = err
				}
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data ItunesResponseModel[ItunesAlbum]
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		return nil, err
	}

	data.Results = slices.Delete(data.Results, 0, 1)

	return &data, nil
}

func (i ItunesAlbumProvider) Name() string {
	return ITUNES_ALBUM_PROVIDER_NAME
}

func (i ItunesAlbumProvider) FetchOutputs(input ItunesArtist) (_ []ItunesAlbum, eOut error) {
	albums, err := fetchAlbums(input.ArtistID)
	if err != nil {
		i.log.Warn(
			"Error fetching album releases",
			slog.Any("err", err),
		)

		return nil, err
	}

	err = i.validate.Struct(albums)
	if err != nil {
		i.log.Warn(
			"Error validating album releases",
			slog.Any("err", err),
		)

		return nil, err
	}

	albums.Results = slices.DeleteFunc(albums.Results, func(output ItunesAlbum) bool {
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

	return albums.Results, nil
}

func (i ItunesAlbumProvider) FromOutputToPersistance(
	inputPersistance model.Inputs,
	output ItunesAlbum,
) (*model.Outputs, error) {
	encodedRaw, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return &model.Outputs{
		ID:            strconv.FormatInt(output.CollectionID, 10),
		InputID:       inputPersistance.ID,
		InputProvider: inputPersistance.Provider,
		Provider:      i.Name(),
		Raw:           encodedRaw,
	}, nil
}
