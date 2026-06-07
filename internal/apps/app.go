package apps

import "context"

type IApp interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
