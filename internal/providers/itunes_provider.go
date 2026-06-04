package providers

type ItunesResponseModel[T any] struct {
	ResultCount int `json:"resultCount" validate:"required,gte=0"`
	Results     []T `json:"results"     validate:"required,dive,required"`
}
