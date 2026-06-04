package jobs

import (
	"context"
	"time"
)

type IJob interface {
	Run(ctx context.Context) error
}

func SleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
		println("hhhhh2")
	}
}
