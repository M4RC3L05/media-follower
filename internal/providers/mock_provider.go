package providers

import (
	"slices"

	"github.com/m4rc3l05/media-follower/.gen/jetdb/model"
)

type CallInfo struct {
	Args []any
}

type MockInputProvider struct {
	Calls       map[string][]CallInfo
	MockReturns map[string][][]any
}

type MockOutputProvider struct {
	Calls       map[string][]CallInfo
	MockReturns map[string][][]any
}

var (
	_ IInputProvider[any]       = MockInputProvider{}
	_ IOutputProvider[any, any] = MockOutputProvider{}
)

func (m MockInputProvider) FromPersistanceToInput(inputPersistance model.Inputs) (*any, error) {
	_, ok := m.Calls["FromPersistanceToInput"]

	if !ok {
		m.Calls["FromPersistanceToInput"] = []CallInfo{}
	}

	m.Calls["FromPersistanceToInput"] = append(
		m.Calls["FromPersistanceToInput"],
		CallInfo{Args: []any{inputPersistance}},
	)

	returns, ok := m.MockReturns["FromPersistanceToInput"]
	if ok {
		if len(returns) <= 0 {
			panic("no more return mocks")
		}

		retValues := returns[0]
		m.MockReturns["FromPersistanceToInput"] = slices.Delete(returns, 0, 1)
		var x *any
		var y error

		if retValues[1] == nil {
			y = nil
		} else {
			y = retValues[1].(error)
		}

		if retValues[0] == nil {
			x = nil
		} else {
			x = any(retValues[0]).(*any)
		}

		return x, y
	} else {
		panic("must mock")
	}
}

func (m MockInputProvider) Name() string {
	_, ok := m.Calls["Name"]

	if !ok {
		m.Calls["Name"] = []CallInfo{}
	}

	m.Calls["Name"] = append(
		m.Calls["Name"],
		CallInfo{Args: []any{}},
	)

	returns, ok := m.MockReturns["Name"]
	if ok {
		if len(returns) <= 0 {
			panic("no more return mocks")
		}

		retValues := returns[0]
		m.MockReturns["Name"] = slices.Delete(returns, 0, 1)

		return retValues[0].(string)
	} else {
		panic("must mock")
	}
}

func (m MockOutputProvider) FetchOutputs(input any) ([]any, error) {
	_, ok := m.Calls["FetchOutputs"]

	if !ok {
		m.Calls["FetchOutputs"] = []CallInfo{}
	}

	m.Calls["FetchOutputs"] = append(
		m.Calls["FetchOutputs"],
		CallInfo{Args: []any{input}},
	)

	returns, ok := m.MockReturns["FetchOutputs"]
	if ok {
		if len(returns) <= 0 {
			panic("no more return mocks")
		}

		retValues := returns[0]
		m.MockReturns["FetchOutputs"] = slices.Delete(returns, 0, 1)

		var x []any
		var y error

		if retValues[1] == nil {
			y = nil
		} else {
			y = retValues[1].(error)
		}

		if retValues[0] == nil {
			x = nil
		} else {
			x = retValues[0].([]any)
		}

		return x, y
	} else {
		panic("must mock")
	}
}

func (m MockOutputProvider) FromOutputToPersistance(
	inputPersistance model.Inputs,
	output any,
) (*model.Outputs, error) {
	_, ok := m.Calls["FromOutputToPersistance"]

	if !ok {
		m.Calls["FromOutputToPersistance"] = []CallInfo{}
	}

	m.Calls["FromOutputToPersistance"] = append(
		m.Calls["FromOutputToPersistance"],
		CallInfo{Args: []any{inputPersistance, output}},
	)

	returns, ok := m.MockReturns["FromOutputToPersistance"]
	if ok {
		if len(returns) <= 0 {
			panic("no more return mocks")
		}

		retValues := returns[0]
		m.MockReturns["FromOutputToPersistance"] = slices.Delete(returns, 0, 1)

		var x *model.Outputs
		var y error

		if retValues[1] == nil {
			y = nil
		} else {
			y = retValues[1].(error)
		}

		if retValues[0] == nil {
			x = nil
		} else {
			x = retValues[0].(*model.Outputs)
		}

		return x, y
	} else {
		panic("must mock")
	}
}

func (m MockOutputProvider) Name() string {
	_, ok := m.Calls["Name"]

	if !ok {
		m.Calls["Name"] = []CallInfo{}
	}

	m.Calls["Name"] = append(
		m.Calls["Name"],
		CallInfo{Args: []any{}},
	)

	returns, ok := m.MockReturns["Name"]
	if ok {
		if len(returns) <= 0 {
			panic("no more return mocks")
		}

		retValues := returns[0]
		m.MockReturns["Name"] = slices.Delete(returns, 0, 1)

		return retValues[0].(string)
	} else {
		panic("must mock")
	}
}
