package database

import (
	"database/sql"
	"fmt"
	"todo-api/internal/config"

	"github.com/go-sql-driver/mysql"
)

func Connected(conf config.Config) (*sql.DB, error) {
	fullURL := conf.DatabaseURL + ":" + conf.DatabasePORT
	cfg := mysql.Config{
		User:      conf.DatabaseUser,
		Passwd:    conf.DatabasePsw,
		Net:       "tcp",
		Addr:      fullURL,
		DBName:    conf.DatabaseName,
		ParseTime: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	fmt.Println("Connected!")

	return db, nil
}
