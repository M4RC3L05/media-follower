package fetchreleases_test

import (
	"context"
	"errors"

	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/common/testing"
	fetchreleases "github.com/m4rc3l05/media-follower/internal/entrypoints/fetch_releases"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FetchReleases", func() {
	Describe("Run()", func() {
		It("should do nothing if no inputs exist", func() {
			db := testing.NewTestDatabase()

			fr := fetchreleases.FetchReleasesEntrypoint{
				DB: db,
				Provider: providers.MockReleaseProvider{
					InterfaceMock: testing.InterfaceMock{
						Calls: map[string][]testing.CallInfo{},
						MockReturns: map[string][][]any{
							"Name": {{providers.ProviderName("Foo")}},
						},
					},
				},
			}

			err := fr.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(fr.Provider).To(testing.HaveInterfaceBeenCalledTimes(1))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("Name", 1))
			Expect(fr.Provider).To(testing.HaveNthMethodBeenCalledWith("Name", 0))
			r := db.DB.QueryRow("select count(*) from releases;")
			var count int
			err = r.Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("should do nothing if no relases for input", func() {
			db := testing.NewTestDatabase()

			fr := fetchreleases.FetchReleasesEntrypoint{
				DB: db,
				Provider: providers.MockReleaseProvider{
					InterfaceMock: testing.InterfaceMock{
						Calls: map[string][]testing.CallInfo{},
						MockReturns: map[string][][]any{
							"Name":          {{providers.ProviderName("Foo")}},
							"FetchReleases": {{[]model.Releases{}, nil}},
						},
					},
				},
			}

			_, err := table.Inputs.
				INSERT(
					table.Inputs.ID,
					table.Inputs.Name,
					table.Inputs.Provider,
					table.Inputs.InternalProviderID,
				).
				VALUES("foo", "bar", "Foo", "foo_internal").
				Exec(db.DB)

			Expect(err).To(BeNil())

			err = fr.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(fr.Provider).To(testing.HaveInterfaceBeenCalledTimes(2))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("Name", 1))
			Expect(fr.Provider).To(testing.HaveNthMethodBeenCalledWith("Name", 0))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(
				fr.Provider,
			).To(testing.HaveNthMethodBeenCalledWith("FetchReleases", 0, model.Inputs{
				ID:                 "foo",
				InternalProviderID: "foo_internal",
				Provider:           "Foo",
				Name:               "bar",
				Description:        nil,
				ImageURL:           nil,
				ExternalLink:       nil,
			}))
			r := db.DB.QueryRow("select count(*) from releases;")
			var count int
			err = r.Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("should not return an error if fetching releases fail", func() {
			db := testing.NewTestDatabase()

			fr := fetchreleases.FetchReleasesEntrypoint{
				DB: db,
				Provider: providers.MockReleaseProvider{
					InterfaceMock: testing.InterfaceMock{
						Calls: map[string][]testing.CallInfo{},
						MockReturns: map[string][][]any{
							"Name":          {{providers.ProviderName("Foo")}},
							"FetchReleases": {{[]model.Releases{}, errors.New("Foo")}},
						},
					},
				},
			}

			_, err := table.Inputs.
				INSERT(
					table.Inputs.ID,
					table.Inputs.Name,
					table.Inputs.Provider,
					table.Inputs.InternalProviderID,
				).
				VALUES("foo", "bar", "Foo", "foo_internal").
				Exec(db.DB)

			Expect(err).To(BeNil())

			err = fr.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(fr.Provider).To(testing.HaveInterfaceBeenCalledTimes(2))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("Name", 1))
			Expect(fr.Provider).To(testing.HaveNthMethodBeenCalledWith("Name", 0))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(
				fr.Provider,
			).To(testing.HaveNthMethodBeenCalledWith("FetchReleases", 0, model.Inputs{
				ID:                 "foo",
				InternalProviderID: "foo_internal",
				Provider:           "Foo",
				Name:               "bar",
				Description:        nil,
				ImageURL:           nil,
				ExternalLink:       nil,
			}))
			r := db.DB.QueryRow("select count(*) from releases;")
			var count int
			err = r.Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("should persist releases", func() {
			db := testing.NewTestDatabase()

			fr := fetchreleases.FetchReleasesEntrypoint{
				DB: db,
				Provider: providers.MockReleaseProvider{
					InterfaceMock: testing.InterfaceMock{
						Calls: map[string][]testing.CallInfo{},
						MockReturns: map[string][][]any{
							"Name": {{providers.ProviderName("Foo")}},
							"FetchReleases": {
								{
									[]model.Releases{
										{
											InternalProviderID: "internal",
											InputID:            "foo",
											Title:              "bar",
											ReleasedAt:         "2026-07-02T19:00:30.040Z",
										},
									},
									nil,
								},
							},
						},
					},
				},
			}

			_, err := table.Inputs.
				INSERT(
					table.Inputs.ID,
					table.Inputs.Name,
					table.Inputs.Provider,
					table.Inputs.InternalProviderID,
				).
				VALUES("foo", "bar", "Foo", "foo_internal").
				Exec(db.DB)

			Expect(err).To(BeNil())

			err = fr.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(fr.Provider).To(testing.HaveInterfaceBeenCalledTimes(2))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("Name", 1))
			Expect(fr.Provider).To(testing.HaveNthMethodBeenCalledWith("Name", 0))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(
				fr.Provider,
			).To(testing.HaveNthMethodBeenCalledWith("FetchReleases", 0, model.Inputs{
				ID:                 "foo",
				InternalProviderID: "foo_internal",
				Provider:           "Foo",
				Name:               "bar",
				Description:        nil,
				ImageURL:           nil,
				ExternalLink:       nil,
			}))
			r := db.DB.QueryRow("select count(*) from releases;")
			var count int
			err = r.Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})

		It("should update existing release", func() {
			db := testing.NewTestDatabase()

			fr := fetchreleases.FetchReleasesEntrypoint{
				DB: db,
				Provider: providers.MockReleaseProvider{
					InterfaceMock: testing.InterfaceMock{
						Calls: map[string][]testing.CallInfo{},
						MockReturns: map[string][][]any{
							"Name": {{providers.ProviderName("Foo")}},
							"FetchReleases": {
								{
									[]model.Releases{
										{
											InternalProviderID: "internal",
											InputID:            "foo",
											Title:              "bar",
											ReleasedAt:         "2025-07-02T19:00:30.040Z",
											Description:        new("foo"),
											ImageURL:           new("https://example.com"),
											ExternalLink:       new("https://bar.com"),
										},
									},
									nil,
								},
							},
						},
					},
				},
			}

			_, err := table.Inputs.
				INSERT(
					table.Inputs.ID,
					table.Inputs.Name,
					table.Inputs.Provider,
					table.Inputs.InternalProviderID,
				).
				VALUES("foo", "bar", "Foo", "foo_internal").
				Exec(db.DB)

			Expect(err).To(BeNil())

			_, err = table.Releases.
				INSERT(
					table.Releases.ID,
					table.Releases.InternalProviderID,
					table.Releases.Title,
					table.Releases.ReleasedAt,
					table.Releases.InputID,
				).
				VALUES("foo", "internal", "title", "2026-07-02T19:00:30.040Z", "foo").
				Exec(db.DB)

			Expect(err).To(BeNil())
			var initRelease model.Releases
			err = table.Releases.SELECT(table.Releases.AllColumns).Query(db.DB, &initRelease)
			Expect(err).To(BeNil())
			Expect(initRelease).To(Equal(model.Releases{
				ID:                 "foo",
				InternalProviderID: "internal",
				InputID:            "foo",
				Title:              "title",
				Description:        nil,
				ImageURL:           nil,
				ExternalLink:       nil,
				ReleasedAt:         "2026-07-02T19:00:30.040Z",
			}))

			err = fr.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(fr.Provider).To(testing.HaveInterfaceBeenCalledTimes(2))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("Name", 1))
			Expect(fr.Provider).To(testing.HaveNthMethodBeenCalledWith("Name", 0))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(
				fr.Provider,
			).To(testing.HaveNthMethodBeenCalledWith("FetchReleases", 0, model.Inputs{
				ID:                 "foo",
				InternalProviderID: "foo_internal",
				Provider:           "Foo",
				Name:               "bar",
				Description:        nil,
				ImageURL:           nil,
				ExternalLink:       nil,
			}))
			r := db.DB.QueryRow("select count(*) from releases;")
			var count int
			err = r.Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
			var updatedRelease model.Releases
			err = table.Releases.SELECT(table.Releases.AllColumns).Query(db.DB, &updatedRelease)
			Expect(err).To(BeNil())
			Expect(updatedRelease).To(Equal(model.Releases{
				ID:                 "foo",
				InternalProviderID: "internal",
				InputID:            "foo",
				Title:              "bar",
				Description:        new("foo"),
				ImageURL:           new("https://example.com"),
				ExternalLink:       new("https://bar.com"),
				ReleasedAt:         "2025-07-02T19:00:30.040Z",
			}))
		})

		It("should do nothing if persist releases fail", func() {
			db := testing.NewTestDatabase()

			fr := fetchreleases.FetchReleasesEntrypoint{
				DB: db,
				Provider: providers.MockReleaseProvider{
					InterfaceMock: testing.InterfaceMock{
						Calls: map[string][]testing.CallInfo{},
						MockReturns: map[string][][]any{
							"Name": {{providers.ProviderName("Foo")}},
							"FetchReleases": {
								{
									[]model.Releases{
										{
											InternalProviderID: "internal",
											InputID:            "foo",
											Title:              "bar",
											ReleasedAt:         "foo",
										},
									},
									nil,
								},
							},
						},
					},
				},
			}

			_, err := table.Inputs.
				INSERT(
					table.Inputs.ID,
					table.Inputs.Name,
					table.Inputs.Provider,
					table.Inputs.InternalProviderID,
				).
				VALUES("foo", "bar", "Foo", "foo_internal").
				Exec(db.DB)

			Expect(err).To(BeNil())

			err = fr.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(fr.Provider).To(testing.HaveInterfaceBeenCalledTimes(2))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("Name", 1))
			Expect(fr.Provider).To(testing.HaveNthMethodBeenCalledWith("Name", 0))
			Expect(fr.Provider).To(testing.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(
				fr.Provider,
			).To(testing.HaveNthMethodBeenCalledWith("FetchReleases", 0, model.Inputs{
				ID:                 "foo",
				InternalProviderID: "foo_internal",
				Provider:           "Foo",
				Name:               "bar",
				Description:        nil,
				ImageURL:           nil,
				ExternalLink:       nil,
			}))
			r := db.DB.QueryRow("select count(*) from releases;")
			var count int
			err = r.Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})
	})
})
