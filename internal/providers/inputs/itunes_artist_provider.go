package inputs

import (
	"github.com/go-playground/validator/v10"
)

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

func (i ItunesArtistProvider) Validate(input ItunesArtist) error {
	return i.validate.Struct(input)
}
