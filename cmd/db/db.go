package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver
	"github.com/leomarzochi/facebooklike/cmd/config"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DBConnection)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
