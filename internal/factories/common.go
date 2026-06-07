package factories

import (
	"context"
	"errors"

	"github.com/m4rc3l05/media-follower/internal/common"
	store "github.com/m4rc3l05/media-follower/internal/storage"
)

func dbFactory(ctx context.Context, cfg *common.Config) (*store.Db, error) {
	db, err := store.New(cfg.Database.Path)
	if err != nil {
		if db != nil {
			if e := db.Close(ctx); e != nil {
				err = errors.Join(err, e)
			}
		}

		return nil, err
	}

	return db, nil
}
