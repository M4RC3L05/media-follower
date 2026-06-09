package providers

import (
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/internal/test"
)

type MockReleaseProvider struct {
	test.InterfaceMock
}

var _ IReleaseProvider[any, any] = MockReleaseProvider{}

func (m MockReleaseProvider) FetchReleases(input any) ([]any, error) {
	test.RecordCall(m.InterfaceMock, "FetchReleases", input)

	var x []any
	var y error
	test.MockReturn(m.InterfaceMock, "FetchReleases", &x, &y)

	return x, y
}

func (m MockReleaseProvider) FromPersistanceToInput(inputPersistance model.Inputs) (*any, error) {
	test.RecordCall(m.InterfaceMock, "FromPersistanceToInput", inputPersistance)

	var x *any
	var y error
	test.MockReturn(m.InterfaceMock, "FromPersistanceToInput", &x, &y)

	return x, y
}

func (m MockReleaseProvider) FromReleaseToPersistance(
	inputPersistance model.Inputs,
	output any,
) (*model.Releases, error) {
	test.RecordCall(m.InterfaceMock, "FromReleaseToPersistance", inputPersistance, output)

	var x *model.Releases
	var y error
	test.MockReturn(m.InterfaceMock, "FromReleaseToPersistance", &x, &y)

	return x, y
}

func (m MockReleaseProvider) Name() string {
	test.RecordCall(m.InterfaceMock, "Name")

	var x string
	test.MockReturn(m.InterfaceMock, "Name", &x)

	return x
}
