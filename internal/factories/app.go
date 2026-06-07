package factories

import (
	"context"
	"fmt"

	"github.com/m4rc3l05/media-follower/internal/common"
	adminserver "github.com/m4rc3l05/media-follower/internal/entrypoints/apps/admin_server"
)

func adminServerAppFactory(ctx context.Context) (*adminserver.AdminServerApp, error) {
	cfg, err := common.NewConfig()
	if err != nil {
		return nil, err
	}

	db, err := dbFactory(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &adminserver.AdminServerApp{Db: db}, nil
}

func appFactory(ctx context.Context, name string, args ...string) (common.IEntrypoint, error) {
	if name == "admin-server" {
		return adminServerAppFactory(ctx)
	}

	return nil, fmt.Errorf("app \"%s\" is not supported", name)
}
