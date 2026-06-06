package jobs_test

import (
	"context"
	"errors"

	"github.com/m4rc3l05/media-follower/.gen/jetdb/model"
	"github.com/m4rc3l05/media-follower/.gen/jetdb/table"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/jobs"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/store"
	"github.com/m4rc3l05/media-follower/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FetchOutputs", func() {
	Describe("Run()", func() {
		It("should do nothing if no inputs exist", func() {
			db := test.NewDB()
			defer db.Close() //nolint:all

			inputProvier := providers.MockInputProvider{
				Calls:       map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{"Name": {{"foo"}}},
			}
			outputProvider := providers.MockOutputProvider{Calls: map[string][]providers.CallInfo{}}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier.Calls).To(HaveLen(1))
			Expect(inputProvier.Calls["Name"]).To(HaveLen(1))
			Expect(outputProvider.Calls).To(BeEmpty())
		})

		It("should do nothing if no inputs exist for provider name", func() {
			db := test.NewDB()
			defer db.Close() //nolint:all

			test.LoadInput(db)

			inputProvier := providers.MockInputProvider{
				Calls:       map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{"Name": {{"biz"}}},
			}
			outputProvider := providers.MockOutputProvider{Calls: map[string][]providers.CallInfo{}}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier.Calls).To(HaveLen(1))
			Expect(inputProvier.Calls["Name"]).To(HaveLen(1))
			Expect(outputProvider.Calls).To(BeEmpty())
		})

		It("should do nothing if `FromPersistanceToInput()` returns error", func() {
			db := test.NewDB()
			defer db.Close() //nolint:all

			test.LoadInput(db)
			callErr := errors.New("foo")
			inputProvier := providers.MockInputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"Name":                   {{"foo"}},
					"FromPersistanceToInput": {{nil, callErr}},
				},
			}
			outputProvider := providers.MockOutputProvider{Calls: map[string][]providers.CallInfo{}}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier.Calls).To(HaveLen(2))
			Expect(inputProvier.Calls["Name"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"]).To(HaveLen(1))
			Expect(outputProvider.Calls).To(BeEmpty())
		})

		It("should do nothing if no outputs are generated", func() {
			db := test.NewDB()
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)

			input := any(struct{}{})
			inputProvier := providers.MockInputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"Name":                   {{"foo"}},
					"FromPersistanceToInput": {{&input, nil}},
				},
			}
			outputProvider := providers.MockOutputProvider{
				Calls:       map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{"FetchOutputs": {{[]any{}, nil}}},
			}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier.Calls).To(HaveLen(2))
			Expect(inputProvier.Calls["Name"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"][0].Args).To(Equal([]any{inputDb}))
			Expect(outputProvider.Calls).To(HaveLen(1))
			Expect(outputProvider.Calls["FetchOutputs"]).To(HaveLen(1))
			Expect(outputProvider.Calls["FetchOutputs"][0].Args).To(Equal([]any{input}))
		})

		It("should do nothing if `FetchOutputs()` returns error", func() {
			db := test.NewDB()
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)
			callErr := errors.New("foo")
			input := any(struct{}{})
			inputProvier := providers.MockInputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"Name":                   {{"foo"}},
					"FromPersistanceToInput": {{&input, nil}},
				},
			}
			outputProvider := providers.MockOutputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"FetchOutputs": {{[]any{}, callErr}},
				},
			}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier.Calls).To(HaveLen(2))
			Expect(inputProvier.Calls["Name"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"][0].Args).To(Equal([]any{inputDb}))
			Expect(outputProvider.Calls).To(HaveLen(1))
			Expect(outputProvider.Calls["FetchOutputs"]).To(HaveLen(1))
			Expect(outputProvider.Calls["FetchOutputs"][0].Args).To(Equal([]any{input}))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(0))
		})

		It("should do nothing if `FromOutputToPersistance` returns error", func() {
			db := test.NewDB()
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)
			callErr := errors.New("foo")
			input := any(struct{}{})
			output := any(struct{}{})
			inputProvier := providers.MockInputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"Name":                   {{"foo"}},
					"FromPersistanceToInput": {{&input, nil}},
				},
			}
			outputProvider := providers.MockOutputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"FetchOutputs":            {{[]any{output}, nil}},
					"FromOutputToPersistance": {{nil, callErr}},
				},
			}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier.Calls).To(HaveLen(2))
			Expect(inputProvier.Calls["Name"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"][0].Args).To(Equal([]any{inputDb}))
			Expect(outputProvider.Calls).To(HaveLen(2))
			Expect(outputProvider.Calls["FetchOutputs"]).To(HaveLen(1))
			Expect(outputProvider.Calls["FetchOutputs"][0].Args).To(Equal([]any{input}))
			Expect(outputProvider.Calls["FromOutputToPersistance"]).To(HaveLen(1))
			Expect(
				outputProvider.Calls["FromOutputToPersistance"][0].Args,
			).To(Equal([]any{inputDb, output}))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(0))
		})

		It("should insert output in db", func() {
			db := test.NewDB()
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)

			input := any(struct{}{})
			output := any(struct{}{})
			outputDb := model.Outputs{
				ID:            "1",
				InputID:       inputDb.ID,
				InputProvider: inputDb.Provider,
				Provider:      "foo-out",
				Raw:           []byte("{}"),
			}
			inputProvier := providers.MockInputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"Name":                   {{"foo"}},
					"FromPersistanceToInput": {{&input, nil}},
				},
			}
			outputProvider := providers.MockOutputProvider{
				Calls: map[string][]providers.CallInfo{},
				MockReturns: map[string][][]any{
					"FetchOutputs":            {{[]any{output}, nil}},
					"FromOutputToPersistance": {{&outputDb, nil}},
				},
			}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)

			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier.Calls).To(HaveLen(2))
			Expect(inputProvier.Calls["Name"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"]).To(HaveLen(1))
			Expect(inputProvier.Calls["FromPersistanceToInput"][0].Args).To(Equal([]any{inputDb}))
			Expect(outputProvider.Calls).To(HaveLen(2))
			Expect(outputProvider.Calls["FetchOutputs"]).To(HaveLen(1))
			Expect(outputProvider.Calls["FetchOutputs"][0].Args).To(Equal([]any{input}))
			Expect(outputProvider.Calls["FromOutputToPersistance"]).To(HaveLen(1))
			Expect(
				outputProvider.Calls["FromOutputToPersistance"][0].Args,
			).To(Equal([]any{inputDb, output}))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0]).To(Equal(outputDb))
		})
	})
})
