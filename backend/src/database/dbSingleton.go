package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

var (
	instance *Database
	once     sync.Once
)

func GetDatabaseInstance() (*Database, error) {
	var err error
	once.Do(func() {
		db, err := sql.Open("sqlite3", "./database/data.db")
		if err != nil {
			return
		}
		instance = &Database{DB: db}
	})
	return instance, err
}
