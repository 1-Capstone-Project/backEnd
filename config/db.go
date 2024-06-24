package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	dsn := "admin:gitmate1234@tcp(gitmate-database.cbimo8eqih5y.ap-northeast-2.rds.amazonaws.com:3306)/gitmate_db?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	fmt.Println("Database connection established")
	return db, nil
}
