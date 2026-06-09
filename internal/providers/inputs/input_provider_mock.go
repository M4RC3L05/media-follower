package inputs

import (
	"github.com/m4rc3l05/media-follower/internal/test"
)

type MockInputProvider[I any] struct {
	test.InterfaceMock
}

func (m MockInputProvider[I]) Validate(input I) error {
	test.RecordCall(m.InterfaceMock, "Validate", input)

	var x error
	test.MockReturn(m.InterfaceMock, "Validate", &x)

	return x
}
