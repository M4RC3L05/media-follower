package jobs

import (
	"context"
)

type IJob interface {
	Run(ctx context.Context) error
}
