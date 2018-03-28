//http://go-database-sql.org/
package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
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
		SELECT sorting_number, date_received, name, building, room, carrier, package_type, is_printed
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
		var s string
		res.Scan(&s, &p.DateReceived, &p.Name, &p.Building, &p.Room, &p.Carrier, &p.PackageType, &p.Printed)
		p.Number = Atosn(s)
		return p, nil
	}

	// Otherwise, return an empty package, and the error indicating that no package was found
	return Package{}, ErrNoPackageFound
}

func GetToBePrinted(building string) ([]Package, error){
	// Prepare getting the number of packages to be printed
	stmt, err := db.Prepare(`
		SELECT COUNT(*)
		FROM Packages
 		WHERE is_printed = 0 AND building = ?`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return nil, err
	}
	defer stmt.Close()

	// Execute getting the number of packages to be printed
	res, err := stmt.Query(building)
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return nil, err
	}
	defer res.Close()
	return nil

	var count int
	if(res.Next()){
		res.Scan(&count)
	}

	// Prepare getting the package info from the database
	stmt, err = db.Prepare(`
		SELECT sorting_number, date_received, name, room, carrier, package_type
		FROM Packages
 		WHERE is_printed = 0 AND building = ?`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return nil, err
	}
	defer stmt.Close()

	// Execute getting the package info from the database
	res, err = stmt.Query(building)
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return nil, err
	}
	defer res.Close()

	//Create array of packages
	toBePrinted := make([]Package, count)
	for i := 0; i < count; i++ {
		if(res.Next()){
			var p Package
			var s string
			res.Scan(&s, &p.DateReceived, &p.Name, &p.Room, &p.Carrier, &p.PackageType)
			p.Number = Atosn(s)
			toBePrinted[i] = p
		}else{
			return toBePrinted, ErrNoPackageFound
		}
	}

	return toBePrinted, nil
}

func UpdatePackage(sortingNumber string, dateReceived time.Time, name string, building string, room string, carrier string, packageType string, isPrinted bool) error {
	stmt, err := db.Prepare(`
		UPDATE Packages
		SET sorting_number = ?, date_received = ?, name = ?, building = ?, room = ?, carrier = ?, package_type = ?, is_printed = ?
		WHERE sorting_number = ?
		`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dateReceived, name, building, room, carrier, packageType, isPrinted, sortingNumber)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	return nil
}

func AddPackage(name string, building string, room string, carrier string, packageType string) (string, error) {
	stmt, err := db.Prepare(`
		INSERT INTO Packages(sorting_number, date_received, name, building, room, carrier, package_type)
		VALUES(?, DATETIME('now','localtime'), ?, ?, ?, ?, ?)`)
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

	_, err = stmt.Exec(sortingNumber.String(), name, building, room, carrier, packageType)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return "", err
	}

	return sortingNumber.String(), nil
}

func Archive(sortingNumber string, signature string) error {

	// Get package information for the archive
	pack, err := GetPackage(sortingNumber)
	if err != nil {
		log.Println("Error occured while getting package:", err)
		return err
	}

	// Archive package -> Put in the Picked_Up table
	stmt, err := db.Prepare("INSERT INTO Picked_Up VALUES(?, ?, ?, ?, ?, ?, ?, DATETIME('now','localtime'), ?)")
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer stmt.Close()
	num, err := stmt.Exec(
		pack.Number.String(),
		pack.DateReceived,
		pack.Name,
		pack.Building,
		pack.Room,
		pack.Carrier,
		pack.PackageType,
		signature,
	)
	
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	rows, _ := num.RowsAffected()
	if rows == 0 {
		return ErrNoPackageFound
	}

	// Remove archived package from current package table
	delstmt, err := db.Prepare("DELETE FROM Packages WHERE sorting_number=?")
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer delstmt.Close()

	_, err = delstmt.Exec(sortingNumber)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	return nil
}
