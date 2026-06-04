package providers

import "github.com/m4rc3l05/media-follower/.gen/jetdb/model"

type EProviders string

type IProvider interface {
	Name() string
}

type IOutputProvider[I any, O any] interface {
	IProvider

	FetchOutputs(input I) ([]O, error)
	FromOutputToPersistance(inputPersistance model.Inputs, output O) (*model.Outputs, error)
}

type IInputProvider[I any] interface {
	IProvider

	// SearchInputs(term string) ([]I, error)
	FromPersistanceToInput(inputPersistance model.Inputs) (*I, error)
}
