//http://go-database-sql.org/
package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"fmt"
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
		SELECT sorting_number, date_received, name, building, room, package_type
		FROM Packages
 		WHERE sorting_number = ?`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return Package{}, err
	}
	defer stmt.Close()

	// Run the query
	res, err := stmt.Query(sortingNumber)
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return Package{}, err
	}
	defer res.Close()

	// If there is a result...
	if res.Next() {

		// Store the found row into a variable p
		var p Package
		res.Scan(&p.Number, &p.DateReceived &p.Name, &p.Building, &p.Room, &p.PackageType)
		fmt.Printf("%s: %#v\n", sortingNumber, p)

		return p, nil
	}

	// Otherwise, return an empty package, and the error indicating that no package was found
	return Package{}, ErrNoPackageFound
}

func AddPackage(name string, building string, room string, packageType string) (string, error) {
	stmt, err := db.Prepare(`
		INSERT INTO Packages(sorting_number, date_received, name, building, room, package_type)
		VALUES(?, DATETIME('now','localtime'), ?, ?, ?, ?)`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return "", err
	}
	defer stmt.Close()

	sortingNumber, err := getNextSortingNumber()
	if err != nil {
		log.Println("Error occured while getting sorting number:", err)
		return "", err
	}

	_, err = stmt.Exec(sortingNumber.String(), name, building, room, packageType)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return "", err
	}

	return sortingNumber.String(), nil
}

func Archive(sortingNumber string) error {

	//Get package information for the archive
	pickedUpPackage = GetPackage(sortingNumber)
	stmt, err := db.Prepare("INSERT INTO Picked_Up VALUES(?, ?, ?, ?, ?, ?, DATETIME('now','localtime'), NULL)")
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer stmt.Close()
	num, err : = stmt.Exec(pickedUpPackage.Number.String(), pickedUpPackage.DateReceived, pickedUpPackage.Name, pickedUpPackage.Building, pickedUpPackage.Room, pickedUpPackage.PackageType)
	
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	rows, _ := num.RowsAffected()
	if rows == 0 {
		return ErrNoPackageFound
	}

	// TODO add a package archive where we store all packages
	stmt, err := db.Prepare("DELETE FROM Packages WHERE sorting_number=?")

	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer stmt.Close()

	num, err := stmt.Exec(sortingNumber)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	rows, _ := num.RowsAffected()
	if rows == 0 {
		return ErrNoPackageFound
	}

	return nil
}
