package providers

import (
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/internal/common/testing"
)

type MockReleaseProvider struct {
	testing.InterfaceMock
}

var _ IReleaseProvider = MockReleaseProvider{}

func (m MockReleaseProvider) FetchReleases(input model.Inputs) ([]model.Releases, error) {
	testing.RecordCall(m.InterfaceMock, "FetchReleases", input)

	var x []model.Releases
	var y error
	testing.MockReturn(m.InterfaceMock, "FetchReleases", &x, &y)

	return x, y
}

func (m MockReleaseProvider) LookupReleaseInput(term string) (*model.Inputs, error) {
	testing.RecordCall(m.InterfaceMock, "LookupReleaseInput", term)

	var x *model.Inputs
	var y error
	testing.MockReturn(m.InterfaceMock, "LookupReleaseInput", &x, &y)

	return x, y
}

func (m MockReleaseProvider) Name() ProviderName {
	testing.RecordCall(m.InterfaceMock, "Name")

	var x ProviderName
	testing.MockReturn(m.InterfaceMock, "Name", &x)

	return x
}
