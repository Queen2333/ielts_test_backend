package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitializeDB 初始化数据库连接池
func InitializeDB(connectionString string) error {
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// GetDB 返回数据库连接
func GetDB() *sql.DB {
	return db
}
