package jobs_test

import (
	"context"
	"errors"

	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/jobs"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/store"
	"github.com/m4rc3l05/media-follower/internal/test"
	"github.com/m4rc3l05/media-follower/internal/test/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FetchOutputs", func() {
	Describe("Run()", func() {
		It("should do nothing if no inputs exist", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls:       map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{"Name": {{"foo"}}},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{Calls: map[string][]test.CallInfo{}},
			}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(0))
		})

		It("should do nothing if no inputs exist for provider name", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			testdata.LoadDBInput(db)

			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls:       map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{"Name": {{"biz"}}},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{Calls: map[string][]test.CallInfo{}},
			}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(0))
		})

		It("should do nothing if `FromPersistanceToInput()` returns error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)
			callErr := errors.New("foo")
			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{nil, callErr}},
					},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{Calls: map[string][]test.CallInfo{}},
			}
			job := jobs.NewFetchOutputsJob(
				inputProvier,
				outputProvider,
				db,
				common.NewLogger("foo"),
			)
			err := job.Run(context.Background())

			Expect(err).To(BeNil())
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				inputProvier,
			).To(test.HaveMethodBeenNthCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(0))
		})

		It("should do nothing if no outputs are generated", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)

			input := any(struct{}{})
			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
					},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls:       map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{"FetchOutputs": {{[]any{}, nil}}},
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
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				inputProvier,
			).To(test.HaveMethodBeenNthCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(1))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FetchOutputs", 1))
			Expect(outputProvider).To(test.HaveMethodBeenNthCalledWith("FetchOutputs", 0, input))
		})

		It("should do nothing if `FetchOutputs()` returns error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)
			callErr := errors.New("foo")
			input := any(struct{}{})
			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
					},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"FetchOutputs": {{[]any{}, callErr}},
					},
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
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				inputProvier,
			).To(test.HaveMethodBeenNthCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(1))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FetchOutputs", 1))
			Expect(outputProvider).To(test.HaveMethodBeenNthCalledWith("FetchOutputs", 0, input))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(0))
		})

		It("should do nothing if `FromOutputToPersistance` returns error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)
			callErr := errors.New("foo")
			input := any(struct{}{})
			output := struct{}{}
			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
					},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"FetchOutputs":            {{[]any{output}, nil}},
						"FromOutputToPersistance": {{nil, callErr}},
					},
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
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				inputProvier,
			).To(test.HaveMethodBeenNthCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FetchOutputs", 1))
			Expect(outputProvider).To(test.HaveMethodBeenNthCalledWith("FetchOutputs", 0, input))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FromOutputToPersistance", 1))
			Expect(
				outputProvider,
			).To(test.HaveMethodBeenNthCalledWith("FromOutputToPersistance", 0, inputDb, output))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

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
			outputDb := model.Outputs{
				ID:            "1",
				InputID:       inputDb.ID,
				InputProvider: inputDb.Provider,
				Provider:      "foo-out",
				Raw:           []byte("{}"),
			}
			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
					},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"FetchOutputs":            {{[]any{output}, nil}},
						"FromOutputToPersistance": {{&outputDb, nil}},
					},
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
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				inputProvier,
			).To(test.HaveMethodBeenNthCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FetchOutputs", 1))
			Expect(outputProvider).To(test.HaveMethodBeenNthCalledWith("FetchOutputs", 0, input))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FromOutputToPersistance", 1))
			Expect(
				outputProvider,
			).To(test.HaveMethodBeenNthCalledWith("FromOutputToPersistance", 0, inputDb, output))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0]).To(Equal(outputDb))
		})

		It("should update existsing output", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			input := any(struct{}{})
			output := struct{}{}
			inputDb := testdata.LoadDBInput(db)
			outputDb := testdata.LoadDBOutput(db, &inputDb)
			updatedOutputDb := model.Outputs{
				ID:            outputDb.ID,
				InputID:       outputDb.InputID,
				InputProvider: outputDb.InputProvider,
				Provider:      outputDb.Provider,
				Raw:           []byte(`{"ok":true}`),
			}
			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
					},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"FetchOutputs":            {{[]any{output}, nil}},
						"FromOutputToPersistance": {{&updatedOutputDb, nil}},
					},
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
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				inputProvier,
			).To(test.HaveMethodBeenNthCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FetchOutputs", 1))
			Expect(outputProvider).To(test.HaveMethodBeenNthCalledWith("FetchOutputs", 0, input))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FromOutputToPersistance", 1))
			Expect(
				outputProvider,
			).To(test.HaveMethodBeenNthCalledWith("FromOutputToPersistance", 0, inputDb, output))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0]).NotTo(Equal(outputDb))
			Expect(outputs[0]).To(Equal(updatedOutputDb))
		})

		It("should insert output in db and ignore the ones that error", func() {
			db := test.NewDB()
			defer db.Close(context.Background()) //nolint:all

			inputDb := testdata.LoadDBInput(db)

			input := any(struct{}{})
			output := struct{}{}
			outputDb1 := model.Outputs{
				ID:            "1",
				InputID:       inputDb.ID,
				InputProvider: inputDb.Provider,
				Provider:      "foo-out",
				Raw:           []byte("{}"),
			}
			outputDb2 := model.Outputs{
				ID:            "1",
				InputID:       inputDb.ID,
				InputProvider: inputDb.Provider,
				Provider:      "foo-out",
				Raw:           []byte("foo"),
			}
			inputProvier := providers.MockInputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"Name":                   {{"foo"}},
						"FromPersistanceToInput": {{&input, nil}},
					},
				},
			}
			outputProvider := providers.MockOutputProvider{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"FetchOutputs":            {{[]any{output, output}, nil}},
						"FromOutputToPersistance": {{&outputDb1, nil}, {&outputDb2, nil}},
					},
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
			Expect(inputProvier).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("Name", 1))
			Expect(inputProvier).To(test.HaveMethodBeenCalledTimes("FromPersistanceToInput", 1))
			Expect(
				inputProvier,
			).To(test.HaveMethodBeenNthCalledWith("FromPersistanceToInput", 0, inputDb))
			Expect(outputProvider).To(test.HaveInterfaceBeenCalledTimes(2))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FetchOutputs", 1))
			Expect(outputProvider).To(test.HaveMethodBeenNthCalledWith("FetchOutputs", 0, input))
			Expect(outputProvider).To(test.HaveMethodBeenCalledTimes("FromOutputToPersistance", 2))
			Expect(
				outputProvider,
			).To(test.HaveMethodBeenNthCalledWith("FromOutputToPersistance", 0, inputDb, output))
			Expect(
				outputProvider,
			).To(test.HaveMethodBeenNthCalledWith("FromOutputToPersistance", 1, inputDb, output))

			var outputs []model.Outputs
			table.Outputs.SELECT(table.Outputs.AllColumns, store.JSONCol(table.Outputs.Raw).AS("outputs.raw")).

				//nolint:all
				Query(db.DB, &outputs)

			Expect(outputs).To(HaveLen(1))
			Expect(outputs[0]).To(Equal(outputDb1))
			Expect(outputs[0]).NotTo(Equal(outputDb2))
		})
	})
})
