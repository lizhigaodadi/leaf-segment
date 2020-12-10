package common

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"
)

func NewMysqlClient() *sql.DB {
	db_host := os.Getenv("MYSQL_HOST")
	db_username := os.Getenv("MYSQL_USERNAME")
	db_password := os.Getenv("MYSQL_PASSWORD")
	db_name := os.Getenv("MYSQL_DB")
	db_conn_maxidlem, _ := strconv.Atoi(os.Getenv("MYSQL_CONN_MAXIDLE"))
	db_conn_maxopen, _ := strconv.Atoi(os.Getenv("MYSQL_CONN_MAXOPEN"))

	connStr := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", db_username, db_password, db_host, db_name)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Printf("init mysql err %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("ping mysql err: %v", err)
	}

	db.SetMaxIdleConns(db_conn_maxidlem)
	db.SetMaxOpenConns(db_conn_maxopen)
	db.SetConnMaxLifetime(5 * time.Minute)
	fmt.Println("init mysql successc")
	return db
}
