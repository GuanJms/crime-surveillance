package data

import (
	"database/sql"
	"errors"
	"time"
)

const dbTimeout = time.Second * 3

var ErrNoConent error

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(conn *sql.DB) *PostgresRepository {
	ErrNoConent = errors.New("no content error")
	return &PostgresRepository{
		Conn: conn,
	}
}
