package outputs_test

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/jarcoal/httpmock"
	"github.com/m4rc3l05/media-follower/internal/providers/inputs"
	"github.com/m4rc3l05/media-follower/internal/providers/outputs"
	"github.com/m4rc3l05/media-follower/internal/test/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ItunesAlbumProvider", func() {
	Describe("FetchOutputs()", Ordered, func() {
		BeforeEach(func() {
			httpmock.Reset()
		})

		It("should return an error if http request fails", func() {
			err := errors.New("mock error")
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewErrorResponder(err),
			)

			Expect(outputs.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(inputs.ItunesArtist{ArtistID: 1}),
			).Error().To(MatchError(err))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return a validation error if http response is invalid", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{"resultCount": -1, "results": []any{}},
				),
			)

			outputs, err := outputs.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(inputs.ItunesArtist{ArtistID: 1})

			Expect(outputs).To(BeNil())
			_, ok := err.(validator.ValidationErrors)
			Expect(ok).To(BeTrue())
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return a validation error if albums are invalid", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 2,
						"results": []map[string]any{
							{"foo": "bar"},
							{"foo": "bar"},
						},
					},
				),
			)

			outputs, err := outputs.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(inputs.ItunesArtist{ArtistID: 1})

			Expect(outputs).To(BeNil())
			_, ok := err.(validator.ValidationErrors)
			Expect(ok).To(BeTrue())
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return a empty list if no albums", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 1,
						"results": []map[string]any{
							{"foo": "bar"},
						},
					},
				),
			)

			albums, err := outputs.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(inputs.ItunesArtist{ArtistID: 1})

			Expect(err).To(BeNil())
			Expect(albums).To(HaveLen(0))
		})

		It("should return albums", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 6,
						"results": []map[string]any{
							{"foo": "bar"},
							testdata.OkItunesAlbumHttpResponse(),
							testdata.BadItunesAlbumNoReleaseHttpResponse(),
							testdata.BadItunesAlbumDJMixHttpResponse(),
							testdata.BadItunesAlbumDJMix2HttpResponse(),
							testdata.BadItunesAlbumCompilationHttpResponse(),
						},
					},
				),
			)

			outputs, _ := outputs.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(inputs.ItunesArtist{ArtistID: 1})
			Expect(outputs).To(HaveLen(5))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})
	})
})
