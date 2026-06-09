package jobs_test

import (
	"context"
	"errors"
	"time"

	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/jobs"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/storage"
	"github.com/m4rc3l05/media-follower/internal/test"
	"github.com/m4rc3l05/media-follower/internal/test/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FetchReleases", func() {
	Describe("Run()", func() {
		It("should do nothing if no inputs exist", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls:       map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{"Name": {{"biz"}}},
				},
			}
			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
		})

		It("should do nothing if no inputs exist for provider name", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			testdata.LoadDBInput(db)

			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls:       map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{"Name": {{"biz"}}},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
		})

		It("should do nothing if `FromPersistanceToInput()` returns error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)
			callErr := errors.New("foo")
			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{nil, callErr}},
					},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
		})

		It("should do nothing if no outputs are generated", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)

			input := any(struct{}{})
			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
						"FetchReleases":          {{[]any{}, nil}},
					},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(3))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(releaseProvider).To(test.HaveNthMethodBeenCalledWith("FetchReleases", 0, input))
		})

		It("should do nothing if `FetchReleases()` returns error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)
			callErr := errors.New("foo")
			input := any(struct{}{})
			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
						"FetchReleases":          {{[]any{}, callErr}},
					},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(3))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(releaseProvider).To(test.HaveNthMethodBeenCalledWith("FetchReleases", 0, input))

			var outputs []model.Releases
			table.Releases.SELECT(table.Releases.AllColumns, storage.JSONCol(table.Releases.Raw).AS("releases.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(0))
		})

		It("should do nothing if `FromReleaseToPersistance` returns error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)
			callErr := errors.New("foo")
			input := any(struct{}{})
			output := struct{}{}
			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                     {{"foo"}},
						"FromPersistanceToInput":   {{&input, nil}},
						"FetchReleases":            {{[]any{output}, nil}},
						"FromReleaseToPersistance": {{nil, callErr}},
					},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(4))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(releaseProvider).To(test.HaveNthMethodBeenCalledWith("FetchReleases", 0, input))
			Expect(
				releaseProvider,
			).To(test.HaveMethodBeenCalledTimes("FromReleaseToPersistance", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromReleaseToPersistance", 0, inputDb, output))

			var outputs []model.Releases
			table.Releases.SELECT(table.Releases.AllColumns, storage.JSONCol(table.Releases.Raw).AS("releases.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(0))
		})

		It("should insert output in db", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)

			input := any(struct{}{})
			output := struct{}{}
			releaseDb := model.Releases{
				ID:            "1",
				InputID:       inputDb.ID,
				InputProvider: inputDb.Provider,
				Provider:      "foo-out",
				ReleasedAt:    time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				Raw:           []byte("{}"),
			}
			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                     {{"foo"}},
						"FromPersistanceToInput":   {{&input, nil}},
						"FetchReleases":            {{[]any{output}, nil}},
						"FromReleaseToPersistance": {{&releaseDb, nil}},
					},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(4))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(releaseProvider).To(test.HaveNthMethodBeenCalledWith("FetchReleases", 0, input))
			Expect(
				releaseProvider,
			).To(test.HaveMethodBeenCalledTimes("FromReleaseToPersistance", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromReleaseToPersistance", 0, inputDb, output))

			var outputs []model.Releases
			table.Releases.SELECT(table.Releases.AllColumns, storage.JSONCol(table.Releases.Raw).AS("releases.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0]).To(Equal(releaseDb))
		})

		It("should update existsing output", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			input := any(struct{}{})
			output := struct{}{}
			inputDb := testdata.LoadDBInput(db)
			releaseDb := testdata.LoadDBRelease(db, &inputDb)
			updatedReleaseDb := model.Releases{
				ID:            releaseDb.ID,
				InputID:       releaseDb.InputID,
				InputProvider: releaseDb.InputProvider,
				Provider:      releaseDb.Provider,
				ReleasedAt:    time.Unix(0, 0).UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				Raw:           []byte(`{"ok":true}`),
			}
			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                     {{"foo"}},
						"FromPersistanceToInput":   {{&input, nil}},
						"FetchReleases":            {{[]any{output}, nil}},
						"FromReleaseToPersistance": {{&updatedReleaseDb, nil}},
					},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(4))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(releaseProvider).To(test.HaveNthMethodBeenCalledWith("FetchReleases", 0, input))
			Expect(
				releaseProvider,
			).To(test.HaveMethodBeenCalledTimes("FromReleaseToPersistance", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromReleaseToPersistance", 0, inputDb, output))

			var outputs []model.Releases
			table.Releases.SELECT(table.Releases.AllColumns, storage.JSONCol(table.Releases.Raw).AS("releases.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0]).NotTo(Equal(releaseDb))
			Expect(outputs[0]).To(Equal(updatedReleaseDb))
		})

		It("should insert output in db and ignore the ones that error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)

			input := any(struct{}{})
			output := struct{}{}
			releaseDb1 := model.Releases{
				ID:            "1",
				InputID:       inputDb.ID,
				InputProvider: inputDb.Provider,
				Provider:      "foo-out",
				ReleasedAt:    time.Unix(0, 0).UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				Raw:           []byte("{}"),
			}
			releaseDb2 := model.Releases{
				ID:            "1",
				InputID:       inputDb.ID,
				InputProvider: inputDb.Provider,
				Provider:      "foo-out",
				Raw:           []byte("foo"),
			}
			releaseProvider := providers.MockReleaseProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                     {{"foo"}},
						"FromPersistanceToInput":   {{&input, nil}},
						"FetchReleases":            {{[]any{output, output}, nil}},
						"FromReleaseToPersistance": {{&releaseDb1, nil}, {&releaseDb2, nil}},
					},
				},
			}

			job := jobs.NewFetchReleasesJob(
				releaseProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(releaseProvider).To(test.HaveInterfaceBeenCalledTimes(4))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(releaseProvider).To(test.HaveMethodBeenCalledTimes("FetchReleases", 1))
			Expect(releaseProvider).To(test.HaveNthMethodBeenCalledWith("FetchReleases", 0, input))
			Expect(
				releaseProvider,
			).To(test.HaveMethodBeenCalledTimes("FromReleaseToPersistance", 2))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromReleaseToPersistance", 0, inputDb, output))
			Expect(
				releaseProvider,
			).To(test.HaveNthMethodBeenCalledWith("FromReleaseToPersistance", 1, inputDb, output))

			var outputs []model.Releases
			table.Releases.SELECT(table.Releases.AllColumns, storage.JSONCol(table.Releases.Raw).AS("releases.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0]).To(Equal(releaseDb1))
			Expect(outputs[0]).NotTo(Equal(releaseDb2))
		})
	})
})
