package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	_ "modernc.org/sqlite"
)

func NewDB() (*sqlx.DB, error) {
	db, err := otelsqlx.Open("sqlite", "file::memory:?cache=shared",
		otelsql.WithAttributes(semconv.DBSystemSqlite))
	if err != nil {
		return nil, err
	}

	return db, nil
}
