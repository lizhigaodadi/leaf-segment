package dao

import "database/sql"

type Dao struct {
	sql *sql.DB
}

func NewDao(sql *sql.DB) *Dao {
	return &Dao{
		sql: sql,
	}
}
