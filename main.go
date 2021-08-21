package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// cmd.Execute()

	db, err := sql.Open("sqlite3", "./gott.db")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM activities")
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
}
