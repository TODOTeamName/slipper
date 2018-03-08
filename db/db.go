//http://go-database-sql.org/
package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB
var ErrNoPackageFound = errors.New("no package found")

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

func GetPackage(sortingNumber string) (Package, error) {

	// Prepare a statement which gets a package
	stmt, err := db.Prepare(`
		SELECT (sorting_number, name, building, room, package_type)
		FROM Packages
 		WHERE sorting_number = ?`)
	defer stmt.Close()
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return Package{}, err
	}

	// Run the query
	res, err := stmt.Query(sortingNumber)
	defer res.Close()

	// If there is a result...
	if res.Next() {

		// Store the found row into a variable p
		var p Package
		res.Scan(&p.SortingNumber, &p.Name, &p.Building, &p.Room, &p.PackageType)

		return p, nil
	}

	// Otherwise, return an empty package, and the error indicating that no package was found
	return Package{}, ErrNoPackageFound
}

func AddPackage(name string, building string, room string, packageType string) error {
	stmt, err := db.Prepare("INSERT INTO Packages(name, building, room, package_type) VALUES(?,?,?,?)")
	defer stmt.Close()
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
	defer stmt.Close()

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
