package data

import (
	"database/sql"
)

func NewRepository(conn *sql.DB) Repository {
	if conn == nil {
		return nil
	} else {
		return NewPostgresRepository(conn)
	}
}

type AuthModel struct {
	Repo Repository
}
