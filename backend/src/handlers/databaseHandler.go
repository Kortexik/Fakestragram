package handlers

import (
	"backend/src/database"
	"database/sql"
	"log"
)

var DB *sql.DB

func init() {
	dbInstance, err := database.GetDatabaseInstance()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = dbInstance.DB

}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
