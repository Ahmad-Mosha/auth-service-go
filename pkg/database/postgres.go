package database

import (
	"time"

	// Import the pgx driver anonymously so it registers itself.
	// We need this so sqlx knows how to talk to PostgreSQL.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Connect establishes a connection to PostgreSQL using sqlx.
func Connect(dsn string) (*sqlx.DB, error) {
	// sqlx.Connect opens a connection AND calls ping under the hood
	// to make sure the database is actually reachable.
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Idiomatic Go: Always configure database connection pooling limits.
	// Without these, Go might open hundreds of database connections,
	// crashing your Postgres database under heavy load.

	// MaxOpenConns sets the maximum number of active database connections.
	db.SetMaxOpenConns(25)

	// MaxIdleConns sets the maximum number of connections kept alive in the idle pool.
	db.SetMaxIdleConns(25)

	// ConnMaxLifetime sets the max time a connection can be reused before it is closed.
	// Helps release stale connections.
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
