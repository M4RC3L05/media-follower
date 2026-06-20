package entrypoints

import "context"

type IEntrypoint interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}
