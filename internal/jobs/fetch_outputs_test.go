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
			defer db.Close() //nolint:all

			test.LoadInput(db)

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
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)
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
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)

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
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)
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
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)
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
			defer db.Close() //nolint:all

			inputDb := test.LoadInput(db)

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
	})
})
