package outputs

import (
	"github.com/m4rc3l05/media-follower/internal/test"
)

type MockOutputProvider[I any, O any] struct {
	test.InterfaceMock
}

func (m MockOutputProvider[I, O]) FetchOutputs(input I) ([]O, error) {
	test.RecordCall(m.InterfaceMock, "FetchOutputs", input)

	var x []O
	var y error
	test.MockReturn(m.InterfaceMock, "FetchOutputs", &x, &y)

	return x, y
}

func (m MockOutputProvider[I, O]) Validate(output O) error {
	test.RecordCall(m.InterfaceMock, "Validate", output)

	var x error
	test.MockReturn(m.InterfaceMock, "Validate", &x)

	return x
}

func (m MockOutputProvider[I, O]) JSONEncode(output O) ([]byte, error) {
	test.RecordCall(m.InterfaceMock, "JSONEncode", output)

	var x []byte
	var y error
	test.MockReturn(m.InterfaceMock, "JSONEncode", &x, &y)

	return x, y
}
