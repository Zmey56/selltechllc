package dbrepo

import (
	"database/sql"
	"github.com/Zmey56/selltechllc/repository"
)

type DBImpl struct {
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB) repository.DB {
	return &DBImpl{
		DB: conn,
	}
}
