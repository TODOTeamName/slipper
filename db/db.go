//http://go-database-sql.org/
package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func Init(path string) {
	con, err := sql.Open("sqlite3", "file:"+path)
	if err != nil {
		log.Fatal(err)
	}
	db = con
}

func Close() {
	db.Close()
}

func AddPackage(name string, building string, room string, packageType string) error {
	stmt, err := db.Prepare("INSERT INTO Packages(name, building, room, package_type) VALUES(?,?,?,?)")
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	_, err = stmt.Exec(name, building, room, packageType)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	return nil
}

func RemovePackage(sortingNumber string) error {

	// TODO add a package archive where we store all packages
	stmt, err := db.Prepare("DELETE FROM Packages WHERE sorting_number=?")
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	_, err = stmt.Exec(sortingNumber)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	return nil
}
