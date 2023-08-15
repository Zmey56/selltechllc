package dbrepo

import (
	"database/sql"
	"github.com/Zmey56/selltechllc/repository"
)

type DBImpl struct {
	DB *sql.DB
}

type testDBImpl struct {
}

func NewPostgresRepo(conn *sql.DB) repository.DB {
	return &DBImpl{
		DB: conn,
	}
}

//func NewTestingsRepo() repository.DB {
//	return testDBImpl{}
//}
