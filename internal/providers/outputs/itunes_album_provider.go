package outputs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/internal/providers/common"
	"github.com/m4rc3l05/media-follower/internal/providers/inputs"
)

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
}

// Compile time check that providers implement interface
var (
	_ IOutputProvider[inputs.ItunesArtist, ItunesAlbum] = ItunesAlbumProvider{}
)

func NewItunesAlbumProvider(validator *validator.Validate) ItunesAlbumProvider {
	return ItunesAlbumProvider{
		validate: validator,
	}
}

func (i ItunesAlbumProvider) fetchAlbums(
	id int64,
) (_ *common.ItunesResponseModel[ItunesAlbum], eOut error) {
	response, err := http.Get(
		fmt.Sprintf(
			"https://itunes.apple.com/lookup?id=%d&entity=album&media=music&sort=recent&limit=60",
			id,
		),
	)

	defer func() {
		if response == nil || response.Body == nil {
			return
		}

		if err := response.Body.Close(); err != nil {
			if eOut != nil {
				eOut = errors.Join(eOut, err)
			} else {
				eOut = err
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

	var data common.ItunesResponseModel[ItunesAlbum]
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Results) > 0 {
		data.Results = slices.Delete(data.Results, 0, 1)
	}

	err = i.validate.Struct(data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (i ItunesAlbumProvider) FetchOutputs(input inputs.ItunesArtist) (_ []ItunesAlbum, eOut error) {
	albums, err := i.fetchAlbums(input.ArtistID)
	if err != nil {
		return nil, err
	}

	return albums.Results, nil
}

func (i ItunesAlbumProvider) Validate(output ItunesAlbum) error {
	return i.validate.Struct(output)
}

func (i ItunesAlbumProvider) JSONEncode(output ItunesAlbum) ([]byte, error) {
	if err := i.validate.Struct(output); err != nil {
		return nil, err
	}

	return json.Marshal(output)
}
