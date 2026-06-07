package providers

import (
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/internal/test"
)

type MockInputProvider struct {
	test.InterfaceMock
}

type MockOutputProvider struct {
	test.InterfaceMock
}

var (
	_ IInputProvider[any]       = MockInputProvider{}
	_ IOutputProvider[any, any] = MockOutputProvider{}
)

func (m MockInputProvider) FromPersistanceToInput(inputPersistance model.Inputs) (*any, error) {
	test.RecordCall(m.InterfaceMock, "FromPersistanceToInput", inputPersistance)

	var x *any
	var y error
	test.MockReturn(m.InterfaceMock, "FromPersistanceToInput", &x, &y)

	return x, y
}

func (m MockInputProvider) Name() string {
	test.RecordCall(m.InterfaceMock, "Name")

	var x string
	test.MockReturn(m.InterfaceMock, "Name", &x)

	return x
}

func (m MockOutputProvider) FetchOutputs(input any) ([]any, error) {
	test.RecordCall(m.InterfaceMock, "FetchOutputs", input)

	var x []any
	var y error
	test.MockReturn(m.InterfaceMock, "FetchOutputs", &x, &y)

	return x, y
}

func (m MockOutputProvider) FromOutputToPersistance(
	inputPersistance model.Inputs,
	output any,
) (*model.Outputs, error) {
	test.RecordCall(m.InterfaceMock, "FromOutputToPersistance", inputPersistance, output)

	var x *model.Outputs
	var y error
	test.MockReturn(m.InterfaceMock, "FromOutputToPersistance", &x, &y)

	return x, y
}

func (m MockOutputProvider) Name() string {
	test.RecordCall(m.InterfaceMock, "Name")

	var x string
	test.MockReturn(m.InterfaceMock, "Name", &x)

	return x
}
