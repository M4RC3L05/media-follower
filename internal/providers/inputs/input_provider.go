package inputs

type IInputProvider[I any] interface {
	Validate(input I) error
}
