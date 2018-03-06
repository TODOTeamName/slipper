//http://go-database-sql.org/
package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
import "log"

var db *sql.DB


func Init(path string) {
	con, err := sql.Open("sqlite3", "file:" + path)
	if err != nil {
		log.Fatal(err)
	}
	db = con
}

func Close() {
	db.Close()
}

func AddPackage(name string, building string, room string, packageType string) {
	stmt, err := db.Prepare("INSERT INTO Packages(name, building, room, package_type) VALUES(?,?,?,?)")
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
	}
	_, err = stmt.Exec(name, building, room, packageType)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
	}
}
