package providers_test

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jarcoal/httpmock"
	"github.com/m4rc3l05/media-follower/.gen/jetdb/model"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ItunesAlbumProvider", func() {
	Describe("Name()", func() {
		It("should return the provider name", func() {
			Expect(
				providers.NewItunesAlbumProvider(validator.New()).Name(),
			).To(Equal(providers.ITUNES_ALBUM_PROVIDER_NAME))
		})
	})

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

			Expect(providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1}),
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

			outputs, err := providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1})

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
						"resultCount": 2,
						"results": []map[string]any{
							{"foo": "bar"},
							{"foo": "bar"},
						},
					},
				),
			)

			outputs, err := providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1})

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
						"resultCount": 1,
						"results": []map[string]any{
							{"foo": "bar"},
						},
					},
				),
			)

			Expect(providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1})).To(HaveLen(0))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should remove from albums all that are compilations", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 1,
						"results": []map[string]any{
							{"foo": "bar"},
							test.OkItunesAlbumMap(),
							test.BadItunesAlbumCompilationMap(),
						},
					},
				),
			)

			outputs, _ := providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1})
			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0].ArtistName).NotTo(Equal("Various Artists"))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should remove from albums all that are dj mixes", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 1,
						"results": []map[string]any{
							{"foo": "bar"},
							test.OkItunesAlbumMap(),
							test.BadItunesAlbumDJMixMap(),
							test.BadItunesAlbumDJMix2Map(),
						},
					},
				),
			)

			outputs, _ := providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1})
			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0].CollectionCensoredName).NotTo(Equal("Various Artists"))
			Expect(outputs[0].CollectionName).NotTo(Equal("Various Artists"))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should remove from albums all that are not released", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 1,
						"results": []map[string]any{
							{"foo": "bar"},
							test.OkItunesAlbumMap(),
							test.BadItunesAlbumNoReleaseMap(),
						},
					},
				),
			)

			outputs, _ := providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1})
			Expect(outputs).To(HaveLen(1))
			Expect(*outputs[0].ReleaseDate).To(Equal("2026-06-05T14:23:45Z"))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should only include valid albums", func() {
			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 1,
						"results": []map[string]any{
							{"foo": "bar"},
							test.OkItunesAlbumMap(),
							test.BadItunesAlbumNoReleaseMap(),
							test.BadItunesAlbumDJMixMap(),
							test.BadItunesAlbumDJMix2Map(),
							test.BadItunesAlbumCompilationMap(),
						},
					},
				),
			)

			outputs, _ := providers.NewItunesAlbumProvider(validator.New()).
				FetchOutputs(providers.ItunesArtist{ArtistID: 1})
			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0].CollectionID).To(Equal(test.OkItunesAlbumMap()["collectionId"]))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})
	})

	Describe("FromOutputToPersistance()", func() {
		It("should return a validation error if output is invalid", func() {
			persistance, err := providers.NewItunesAlbumProvider(validator.New()).
				FromOutputToPersistance(
					model.Inputs{ID: "foo", Provider: "bar"},
					providers.ItunesAlbum{},
				)

			Expect(persistance).To(BeNil())
			_, ok := err.(validator.ValidationErrors)
			Expect(ok).To(BeTrue())
		})
		It("should return a persistance from output", func() {
			output := test.OkItunesAlbumStruct()
			persistance, err := providers.NewItunesAlbumProvider(validator.New()).
				FromOutputToPersistance(
					model.Inputs{ID: "foo", Provider: "bar"},
					output,
				)

			Expect(err).To(BeNil())
			Expect(persistance).NotTo(BeNil())
			Expect(persistance.ID).To(Equal(strconv.FormatInt(output.CollectionID, 10)))
			Expect(persistance.InputID).To(Equal("foo"))
			Expect(persistance.InputProvider).To(Equal("bar"))
			Expect(persistance.Provider).To(Equal(providers.ITUNES_ALBUM_PROVIDER_NAME))
			Expect(persistance.Raw).NotTo(BeNil())
		})
	})
})
