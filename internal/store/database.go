package store

import (
	"database/sql"
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	qb "github.com/go-jet/jet/v2/sqlite"
	_ "modernc.org/sqlite"
)

type Db struct {
	DB *sql.DB
}

func New() (*Db, error) {
	pp, err := url.JoinPath("./data/", "app.db")
	if err != nil {
		return nil, err
	}

	ppabs, err := filepath.Abs(pp)
	if err != nil {
		return nil, err
	}

	dbUrl := url.URL{Scheme: "file", Path: ppabs, RawQuery: strings.Join([]string{
		"_pragma=journal_mode(WAL)",
		"_pragma=busy_timeout(30000)",
		"_pragma=foreign_keys(ON)",
		"_pragma=synchronous(NORMAL)",
		"_pragma=temp_store(MEMORY)",
	}, "&")}
	db, err := sql.Open("sqlite", dbUrl.String())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)

	return &Db{DB: db}, nil
}

func (db *Db) Close() error {
	var errAgg []error

	if db.DB == nil {
		return nil
	}

	if _, err := db.DB.Exec("pragma optimise"); err != nil {
		errAgg = append(errAgg, err)
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
