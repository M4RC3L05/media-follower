package common

import "context"

type IEntrypoint interface {
	Run(ctx context.Context) error
	Close(ctx context.Context) error
}
