package outputs

type IOutputProvider[I any, O any] interface {
	FetchOutputs(input I) ([]O, error)
	Validate(output O) error
	JSONEncode(output O) ([]byte, error)
}
