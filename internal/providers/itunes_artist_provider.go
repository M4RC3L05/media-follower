package providers

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
)

const ITUNES_ARTIST_PROVIDER_NAME = "ITUNES_ARTIST_PROVIDER"

type ItunesArtist struct {
	WrapperType      string `json:"wrapperType"           validate:"required,eq=artist"`
	ArtistType       string `json:"artistType"            validate:"required,eq=Artist"`
	ArtistName       string `json:"artistName"            validate:"required"`
	ArtistLinkURL    string `json:"artistLinkUrl"         validate:"required,url"`
	ArtistID         int64  `json:"artistId"              validate:"required,number"`
	AmgArtistID      *int64 `json:"amgArtistId,omitempty" validate:"omitempty,number"`
	PrimaryGenreName string `json:"primaryGenreName"      validate:"required"`
	PrimaryGenreID   int64  `json:"primaryGenreId"        validate:"required,number"`
}

type ItunesArtistProvider struct {
	validate *validator.Validate
}

// Compile time check that providers implement interface
var (
	_ IInputProvider[ItunesArtist] = ItunesArtistProvider{}
)

func NewItunesArtistProvider(validator *validator.Validate) ItunesArtistProvider {
	return ItunesArtistProvider{
		validate: validator,
	}
}

func (i ItunesArtistProvider) Name() string {
	return ITUNES_ARTIST_PROVIDER_NAME
}

func (i ItunesArtistProvider) FromPersistanceToInput(
	inputPersistance model.Inputs,
) (*ItunesArtist, error) {
	var itunesArtist ItunesArtist
	err := json.Unmarshal(inputPersistance.Raw, &itunesArtist)
	if err != nil {
		return nil, err
	}

	err = i.validate.Struct(itunesArtist)
	if err != nil {
		return nil, err
	}

	return &itunesArtist, nil
}
