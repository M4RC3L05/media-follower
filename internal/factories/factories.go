package factories

import (
	"context"
	"fmt"

	"github.com/m4rc3l05/media-follower/internal/common"
)

func EntrypointFactory(
	ctx context.Context,
	t string,
	name string,
	args ...string,
) (common.IEntrypoint, error) {
	if t == "app" {
		return appFactory(ctx, name, args...)
	}

	if t == "job" {
		return jobFactory(ctx, name, args...)
	}

	return nil, fmt.Errorf("entrypoint type \"%s\" is not valid", t)
}
