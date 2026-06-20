package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/m4rc3l05/media-follower/internal/common/logging"
	_ "modernc.org/sqlite"
)

var log = logging.New("database")

//go:embed schema.sql
var DBSchema []byte

type Db struct {
	DB *sql.DB
}

func New(path string) (*Db, error) {
	queryParams := strings.Join([]string{
		"_pragma=journal_mode(WAL)",
		"_pragma=busy_timeout(30000)",
		"_pragma=foreign_keys(ON)",
		"_pragma=synchronous(NORMAL)",
		"_pragma=temp_store(MEMORY)",
		"_pragma=optimize(0x10002)",
	}, "&")

	dbUrl := fmt.Sprintf("file:%s?%s", path, queryParams)

	db, err := sql.Open("sqlite", dbUrl)
	if err != nil {
		return nil, err
	}

	nCons := max(runtime.NumCPU(), 1)

	log.Info(fmt.Sprintf("Starting db with %d cons", nCons))

	db.SetMaxOpenConns(nCons)
	db.SetMaxIdleConns(nCons)

	return &Db{DB: db}, nil
}

func (db *Db) Close(ctx context.Context) error {
	var errAgg []error

	if db.DB == nil {
		return nil
	}

	if !testing.Testing() {
		if _, err := db.DB.ExecContext(ctx, "pragma optimize"); err != nil {
			errAgg = append(errAgg, err)
		}
	}

	if err := db.DB.Close(); err != nil {
		errAgg = append(errAgg, err)
	}

	if len(errAgg) <= 0 {
		return nil
	}

	return errors.Join(errAgg...)
}
