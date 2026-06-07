package providers_test

import (
	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/internal/providers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ItunesArtistProvider", func() {
	Describe("Name()", func() {
		It("should return the provider name", func() {
			Expect(
				providers.NewItunesArtistProvider(validator.New()).Name(),
			).To(Equal(providers.ITUNES_ARTIST_PROVIDER_NAME))
		})
	})

	Describe("FromPersistanceToInput()", func() {
		It("should return error if it is unable to parse persitance raw json", func() {
			Expect(
				providers.NewItunesArtistProvider(validator.New()).
					FromPersistanceToInput(model.Inputs{}),
			).Error().Should(MatchError("unexpected end of JSON input"))
		})

		It("should return validation error if persistance is invalid", func() {
			input, err := providers.NewItunesArtistProvider(validator.New()).
				FromPersistanceToInput(model.Inputs{Raw: []byte("{}")})

			Expect(input).To(BeNil())
			_, ok := err.(validator.ValidationErrors)
			Expect(ok).To(BeTrue())
		})

		It("should return a input from persistance", func() {
			input, _ := providers.NewItunesArtistProvider(validator.New()).
				FromPersistanceToInput(model.Inputs{Raw: []byte(`{"wrapperType":"artist","artistType":"Artist","artistName":"The Example Band","artistLinkUrl":"https://example.com/artist/123","artistId":123456789,"amgArtistId":987654321,"primaryGenreName":"Rock","primaryGenreId":21}`)})

			Expect(input.WrapperType).To(Equal("artist"))
			Expect(input.ArtistType).To(Equal("Artist"))
			Expect(input.ArtistName).To(Equal("The Example Band"))
			Expect(input.ArtistLinkURL).To(Equal("https://example.com/artist/123"))
			Expect(input.ArtistID).To(Equal(int64(123456789)))
			Expect(input.AmgArtistID).ToNot(BeNil())
			Expect(*input.AmgArtistID).To(Equal(int64(987654321)))
			Expect(input.PrimaryGenreName).To(Equal("Rock"))
			Expect(input.PrimaryGenreID).To(Equal(int64(21)))
		})
	})
})
