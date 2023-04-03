package postgre

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Conn, returns connection instance for postgre.
func Conn() *pgxpool.Pool {
	return conn
}

// InsertQuery, inserts a new record into the database.
func InsertQuery(query string, args ...interface{}) (pgconn.CommandTag, error) {
	resp, err := conn.Exec(context.Background(), query, args...)
	return resp, err
}

// SelectQuery, selects records from the database.
func SelectQuery(query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := conn.Query(context.Background(), query, args...)
	return rows, err
}
