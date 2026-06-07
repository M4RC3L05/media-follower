package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	qb "github.com/go-jet/jet/v2/sqlite"
	_ "modernc.org/sqlite"
)

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
	}, "&")

	dbUrl := fmt.Sprintf("file:%s?%s", path, queryParams)

	db, err := sql.Open("sqlite", dbUrl)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)

	return &Db{DB: db}, nil
}

func (db *Db) Close(ctx context.Context) error {
	var errAgg []error

	if db.DB == nil {
		return nil
	}

	if !testing.Testing() {
		if _, err := db.DB.ExecContext(ctx, "pragma optimise"); err != nil {
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

func JSONCol(col qb.Column) qb.Expression {
	return qb.Func("json", col)
}

func JSONB(data []byte) qb.BlobExpression {
	return qb.BlobExp(qb.Func("jsonb", qb.Blob(data)))
}
