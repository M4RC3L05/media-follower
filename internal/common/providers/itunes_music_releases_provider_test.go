package providers_test

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/jarcoal/httpmock"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/common/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ItunesMusicReleasesProvider", func() {
	Describe("Name()", func() {
		It("should return provider name", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			Expect(p.Name()).To(Equal(providers.ITUNES_MUSIC_RELEASES_PROVIDER))
		})
	})

	Describe("LookupReleaseInput()", func() {
		BeforeEach(func() {
			httpmock.Activate(GinkgoTB())
		})

		AfterEach(func() {
			httpmock.Reset()
		})

		It("should return error if input lookup request fails", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			err := errors.New("foo")

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewErrorResponder(err),
			)

			res, err := p.LookupReleaseInput("foo")

			Expect(res).To(BeNil())
			Expect(err).To(MatchErrorStrictly(err))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return validation error if data is invalid", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": -1,
						"results":     []any{struct{ artistLinkUrl string }{artistLinkUrl: "foo"}},
					},
				),
			)

			res, err := p.LookupReleaseInput("foo")

			Expect(res).To(BeNil())
			_, ok := err.(validator.ValidationErrors)
			Expect(ok).To(BeTrue())
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return error if no items were returned", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{"resultCount": 0, "results": []any{}},
				),
			)

			res, err := p.LookupReleaseInput("foo")

			Expect(res).To(BeNil())
			Expect(err).To(MatchError("search term \"foo\" did not return any input"))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return error fetching input image fails", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{"resultCount": 0, "results": []providers.ItunesLookupArtist{
						{
							WrapperType:      "artist",
							ArtistType:       "Artist",
							ArtistLinkUrl:    "https://example.com",
							ArtistId:         1,
							ArtistName:       "foo",
							PrimaryGenreId:   2,
							PrimaryGenreName: "bar",
						},
					}},
				),
			)

			e := errors.New("foo")
			httpmock.RegisterResponder(
				"GET",
				"https://example.com",
				httpmock.NewErrorResponder(e),
			)

			res, err := p.LookupReleaseInput("foo")

			Expect(res).To(BeNil())
			Expect(err).To(MatchErrorStrictly(e))
			Expect(httpmock.GetTotalCallCount()).To(Equal(2))
		})

		It("should return a input", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{"resultCount": 0, "results": []providers.ItunesLookupArtist{
						{
							WrapperType:      "artist",
							ArtistType:       "Artist",
							ArtistName:       "foo",
							ArtistLinkUrl:    "https://example.com",
							ArtistId:         2,
							PrimaryGenreId:   3,
							PrimaryGenreName: "biz",
						},
					}},
				),
			)

			httpmock.RegisterResponder(
				"GET",
				"https://example.com",
				httpmock.NewStringResponder(200, `
          <!doctype html>
          <html dir="ltr" lang="en-US">

          <head>
            <meta property="og:title" content="Roosevelt on Apple Music" />
            <meta property="og:image"
              content="https://is1-ssl.mzstatic.com/image/thumb/AMCArtistImages116/v4/b0/e2/d0/b0e2d047-19b5-87c3-8f8e-fd60690db61f/ae3cadda-736a-4e96-b4b9-68f642a94439_file_cropped.png/1200x630cw.png" />
            <meta property="og:image:alt" content="Roosevelt on Apple Music" />
            <meta property="og:image:width" content="1200" />
            <meta property="og:image:height" content="630" />
            <meta property="og:image:type" content="image/png" />
          </head>

          <body></body>

          </html>
        `),
			)

			res, err := p.LookupReleaseInput("foo")

			Expect(res.InternalProviderID).To(Equal("2"))
			Expect(res.Provider).To(Equal(string(p.Name())))
			Expect(res.Name).To(Equal("foo"))
			Expect(
				*res.ImageURL,
			).To(Equal("https://is1-ssl.mzstatic.com/image/thumb/AMCArtistImages116/v4/b0/e2/d0/b0e2d047-19b5-87c3-8f8e-fd60690db61f/ae3cadda-736a-4e96-b4b9-68f642a94439_file_cropped.png/300x300.png"))
			Expect(*res.ExternalLink).To(Equal("https://example.com"))

			Expect(err).To(BeNil())
			Expect(httpmock.GetTotalCallCount()).To(Equal(2))
		})

		It("should return a input without image if it cannot get teh artist image url", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{"resultCount": 0, "results": []providers.ItunesLookupArtist{
						{
							WrapperType:      "artist",
							ArtistType:       "Artist",
							ArtistName:       "foo",
							ArtistLinkUrl:    "https://example.com",
							ArtistId:         2,
							PrimaryGenreId:   3,
							PrimaryGenreName: "biz",
						},
					}},
				),
			)

			httpmock.RegisterResponder(
				"GET",
				"https://example.com",
				httpmock.NewStringResponder(200, `
          <!doctype html>
          <html dir="ltr" lang="en-US">

          <head>
            <meta property="og:title" content="Roosevelt on Apple Music" />
            <meta property="og:image:alt" content="Roosevelt on Apple Music" />
            <meta property="og:image:width" content="1200" />
            <meta property="og:image:height" content="630" />
            <meta property="og:image:type" content="image/png" />
          </head>

          <body></body>

          </html>
        `),
			)

			res, err := p.LookupReleaseInput("foo")

			Expect(res.InternalProviderID).To(Equal("2"))
			Expect(res.Provider).To(Equal(string(p.Name())))
			Expect(res.Name).To(Equal("foo"))
			Expect(
				res.ImageURL,
			).To(BeNil())
			Expect(*res.ExternalLink).To(Equal("https://example.com"))

			Expect(err).To(BeNil())
			Expect(httpmock.GetTotalCallCount()).To(Equal(2))
		})
	})

	Describe("FetchReleases()", func() {
		BeforeEach(func() {
			httpmock.Activate(GinkgoTB())
		})

		AfterEach(func() {
			httpmock.Reset()
		})

		It("should return error if input lookup request fails", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			err := errors.New("foo")

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewErrorResponder(err),
			)

			res, err := p.FetchReleases(model.Inputs{})

			Expect(res).To(BeNil())
			Expect(err).To(MatchErrorStrictly(err))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return validation error if data is invalid", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": -1,
						"results": []any{
							struct{}{},
							struct{ collectionViewURL string }{collectionViewURL: "foo"},
						},
					},
				),
			)

			res, err := p.FetchReleases(model.Inputs{})

			Expect(res).To(BeNil())
			_, ok := err.(validator.ValidationErrors)
			Expect(ok).To(BeTrue())
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return ok if no items are returned", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 0,
						"results": []any{
							struct{}{},
						},
					},
				),
			)

			res, err := p.FetchReleases(model.Inputs{})

			Expect(err).To(BeNil())
			Expect(res).To(Equal([]model.Releases{}))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should normalize release date to utc", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 0,
						"results": []any{
							struct{}{},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               12345,
								"artistName":             "Example Artist",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "Example Album",
								"collectionExplicitness": "cleaned",
								"collectionId":           67890,
								"collectionName":         "Example Album",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            "2026-06-22T08:34:56.789-04:00",
								"trackCount":             12,
								"wrapperType":            "collection",
							},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               12345,
								"artistName":             "Example Artist",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "Example Album",
								"collectionExplicitness": "cleaned",
								"collectionId":           67890,
								"collectionName":         "Example Album",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            "2026-06-22T08:34:56.789Z",
								"trackCount":             12,
								"wrapperType":            "collection",
							},
						},
					},
				),
			)

			res, err := p.FetchReleases(model.Inputs{})

			Expect(err).To(BeNil())
			Expect(res[0].ReleasedAt).To(Equal("2026-06-22T12:34:56.789Z"))
			Expect(res[1].ReleasedAt).To(Equal("2026-06-22T08:34:56.789Z"))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should ignore invalid releases", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 0,
						"results": []any{
							struct{}{},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               1,
								"artistName":             "Example Artist",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "Example Album",
								"collectionExplicitness": "cleaned",
								"collectionId":           1,
								"collectionName":         "Example Album",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            nil,
								"trackCount":             12,
								"wrapperType":            "collection",
							},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               2,
								"artistName":             "Varios Artists",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "Example Album",
								"collectionExplicitness": "cleaned",
								"collectionId":           2,
								"collectionName":         "Example Album",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            "2026-06-22T08:34:56.789Z",
								"trackCount":             12,
								"wrapperType":            "collection",
							},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               3,
								"artistName":             "Exmaple Artist",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "DJ Mix",
								"collectionExplicitness": "cleaned",
								"collectionId":           3,
								"collectionName":         "Example Album",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            "2026-06-22T08:34:56.789Z",
								"trackCount":             12,
								"wrapperType":            "collection",
							},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               4,
								"artistName":             "Exmaple Artist",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "Example Album",
								"collectionExplicitness": "cleaned",
								"collectionId":           4,
								"collectionName":         "DJ Mix",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            "2026-06-22T08:34:56.789Z",
								"trackCount":             12,
								"wrapperType":            "collection",
							},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               5,
								"artistName":             "Example Artist",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "Example Album",
								"collectionExplicitness": "cleaned",
								"collectionId":           5,
								"collectionName":         "Example Album",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            "2026-06-22T08:34:56.789Z",
								"trackCount":             12,
								"wrapperType":            "collection",
							},
						},
					},
				),
			)

			res, err := p.FetchReleases(model.Inputs{})

			Expect(err).To(BeNil())
			Expect(res).To(HaveLen(1))
			Expect(res[0].InternalProviderID).To(Equal("5"))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})

		It("should return releases", func() {
			p := providers.ItunesMusicReleaseProvider{
				Validator: utils.NewValidator(),
				Conform:   utils.NewModifier(),
			}

			httpmock.RegisterResponder(
				"GET",
				"https://itunes.apple.com/lookup",
				httpmock.NewJsonResponderOrPanic(
					200,
					map[string]any{
						"resultCount": 0,
						"results": []any{
							struct{}{},
							map[string]any{
								"amgArtistId":            nil,
								"artistId":               5,
								"artistName":             "Example Artist",
								"artistViewUrl":          "https://example.com/artist/12345",
								"artworkUrl100":          "https://example.com/artwork/100x100.jpg",
								"artworkUrl60":           "https://example.com/artwork/60x60.jpg",
								"collectionCensoredName": "Example Album",
								"collectionExplicitness": "cleaned",
								"collectionId":           5,
								"collectionName":         "Example Album",
								"collectionPrice":        9.99,
								"collectionType":         "Album",
								"collectionViewUrl":      "https://example.com/album/67890",
								"contentAdvisoryRating":  "Teen",
								"copyright":              "© 2026 Example",
								"country":                "US",
								"currency":               "USD",
								"primaryGenreName":       "Pop",
								"releaseDate":            "2026-06-22T08:34:56.789Z",
								"trackCount":             12,
								"wrapperType":            "collection",
							},
						},
					},
				),
			)

			res, err := p.FetchReleases(model.Inputs{ID: "1"})

			Expect(err).To(BeNil())
			Expect(res).To(Equal([]model.Releases{
				{
					ID:                 "",
					InternalProviderID: "5",
					InputID:            "1",
					Title:              "Example Artist - Example Album",
					Description:        nil,
					ImageURL:           new("https://example.com/artwork/300x300bb.jpg"),
					ExternalLink:       new("https://example.com/album/67890"),
					ReleasedAt:         "2026-06-22T08:34:56.789Z",
				},
			}))
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		})
	})
})
