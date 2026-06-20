package utils

import (
	"context"
	"reflect"
	"slices"
	"strconv"
	"time"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}

	return nil
}

type CustomBinder struct {
	Conform   *mold.Transformer
	DefBinder *echo.DefaultBinder
}

func (cv *CustomBinder) Bind(c *echo.Context, i any) error {
	if err := cv.DefBinder.Bind(c, i); err != nil {
		return err
	}

	if err := cv.Conform.Struct(context.Background(), i); err != nil {
		return err
	}

	return nil
}

func NewModifier() *mold.Transformer {
	m := modifiers.New()

	m.Register("min", func(ctx context.Context, fl mold.FieldLevel) error {
		v, err := strconv.ParseInt(fl.Param(), 10, 64)
		if err != nil {
			return err
		}

		if fl.Field().Int() < v {
			fl.Field().SetInt(v)
		}

		return nil
	})

	m.Register("toisoutc", func(ctx context.Context, fl mold.FieldLevel) error {
		if fl.Field().Kind() == reflect.Pointer && fl.Field().IsNil() {
			return nil
		}

		t, err := time.Parse(time.RFC3339Nano, fl.Field().String())
		if err != nil {
			return err
		}

		fl.Field().SetString(t.UTC().Format("2006-01-02T15:04:05.000Z07:00"))

		return nil
	})

	return m
}

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	err := v.RegisterValidation("providerName", func(fl validator.FieldLevel) bool {
		i := slices.Index(providers.PROVIDERS, providers.ProviderName(fl.Field().String()))

		return i != -1
	})
	if err != nil {
		panic(err)
	}

	return v
}
