package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func StartMYSQL() {
	db, err := sql.Open("mysql", "root:chatappdb@tcp(localhost:3306)/chatapp?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	DB = db
}
