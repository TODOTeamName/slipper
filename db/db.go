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

// Initializes the database.
//
// path: the path to the sqlite3 file
func Init(path string) {
	con, err := sql.Open("sqlite3", "file:"+path)
	if err != nil {
		log.Fatal(err)
	}
	db = con
}

// Closes the database.
func Close() {
	db.Close()
}

// Gets package information of a given sorting number and building.
func GetPackage(sortingNumber string, building string) (Package, error) {

	// Prepare a statement which gets a package
	stmt, err := db.Prepare(`
		SELECT sorting_number, date_received, name, building, room, carrier, package_type, is_printed
		FROM Packages
 		WHERE sorting_number = ? AND building = ?`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return Package{}, err
	}
	defer stmt.Close()

	// Run the query
	res, err := stmt.Query(sortingNumber, building)
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

// Gets which packages need to be printed for a building.
func GetToBePrinted(building string) ([]Package, error) {
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

	var count int
	if res.Next() {
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
		if res.Next() {
			var p Package
			var s string
			res.Scan(&s, &p.DateReceived, &p.Name, &p.Room, &p.Carrier, &p.PackageType)
			p.Number = Atosn(s)
			toBePrinted[i] = p
		} else {
			return toBePrinted, ErrNoPackageFound
		}
	}

	return toBePrinted, nil
}

// Mark a building's packages as printed.
func MarkPrinted(building string) error{
	stmt, err := db.Prepare(`
		UPDATE Packages
		SET is_printed = 1
		WHERE is_printed = 0 AND building = ?
		`)
	if err != nil {
		log.Println("Error occured while preparing statement")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(building)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	return nil
}

// Updates a package's info.
func UpdatePackage(sortingNumber string, name string, building string, room string, carrier string, packageType string, isPrinted int) error {
	log.Println(sortingNumber, building, name)
	stmt, err := db.Prepare(`
		UPDATE Packages
		SET name = ?, room = ?, carrier = ?, package_type = ?, is_printed = ?
		WHERE sorting_number = ? and building = ?
		`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, room, carrier, packageType, isPrinted, sortingNumber, building)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	return nil
}

// Adds a package to the database.
func AddPackage(name string, building string, room string, carrier string, packageType string) (string, error) {
	stmt, err := db.Prepare(`
		INSERT INTO Packages(sorting_number, date_received, name, building, room, carrier, package_type)
		VALUES(?, DATETIME('now','localtime'), ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return "", err
	}
	defer stmt.Close()

	sortingNumber, err := getNextSortingNumber(building)
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

// Archives a package in the database. Moves a package from
// the Packages table to the Picked_Up table.
func Archive(sortingNumber string, building string, signature string) error {

	// Get package information for the archive
	pack, err := GetPackage(sortingNumber, building)
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
	delstmt, err := db.Prepare("DELETE FROM Packages WHERE sorting_number=? AND building=?")
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer delstmt.Close()

	_, err = delstmt.Exec(sortingNumber, building)
	if err != nil {
		log.Println("Error occured while executing statement:", err)
		return err
	}

	return nil
}

// Gets all archived packages for a given person in a room.
func CheckArchive(name string, room string, building string) ([]Package, error){
	// Prepare getting the number of packages to be printed
	stmt, err := db.Prepare(`
		SELECT COUNT(*)
		FROM Picked_Up
 		WHERE (name = ? OR room = ?) AND building = ?`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return nil, err
	}
	defer stmt.Close()

	// Execute getting the number of packages to be printed
	res, err := stmt.Query(name, room, building)
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return nil, err
	}
	defer res.Close()

	var count int
	if res.Next() {
		res.Scan(&count)
	}

	// Prepare getting the package info from the database
	stmt, err = db.Prepare(`
		SELECT date_received, name, room, carrier, package_type, date_picked_up, signature
		FROM Picked_Up
 		WHERE (name = ? OR room = ?) AND building = ?`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return nil, err
	}
	defer stmt.Close()

	// Execute getting the package info from the database
	res, err = stmt.Query(name, room, building)
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return nil, err
	}
	defer res.Close()

	//Create array of packages
	fromArchive := make([]Package, count)
	for i := 0; i < count; i++ {
		if res.Next() {
			var p Package
			res.Scan(&p.DateReceived, &p.Name, &p.Room, &p.Carrier, &p.PackageType, &p.DatePickedUp, &p.Signature)
			fromArchive[i] = p
		} else {
			return fromArchive, ErrNoPackageFound
		}
	}

	return fromArchive, nil
}

// Gets the password hash for the login of a building
func GetPassword(building string) (string, error){
	// Prepare a statement which gets a package
	stmt, err := db.Prepare(`
		SELECT password
		FROM Users
 		WHERE building = ?`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return "", err
	}
	defer stmt.Close()

	// Run the query
	res, err := stmt.Query(building)
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return "", err
	}
	defer res.Close()

	// If there is a result...
	if res.Next() {

		// Store the found row into a variable p
		var pass string
		res.Scan(&pass)
		return pass, nil
	}

	return "", err
}

// Cleans the archive so that the database doesn't get clogged up.
func CleanArchive() error{
	// Prepare getting the number of packages to be printed
	stmt, err := db.Prepare(`
		SELECT COUNT(*)
		FROM Picked_Up`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer stmt.Close()

	// Execute getting the number of packages to be printed
	res, err := stmt.Query()
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return err
	}
	defer res.Close()

	var count int
	if res.Next() {
		res.Scan(&count)
	}

	// Prepare getting the package info from the database
	stmt, err = db.Prepare(`
		SELECT date_received, date_picked_up
		FROM Picked_Up`)
	if err != nil {
		log.Println("Error occured while preparing statement:", err)
		return err
	}
	defer stmt.Close()

	// Execute getting the package info from the database
	res, err = stmt.Query()
	if err != nil {
		log.Println("Error occured while executing query:", err)
		return err
	}
	defer res.Close()

	now := time.Now()
	// Iterate over all the packages in Picked_Up
	for i := 0; i < count; i++ {
		if res.Next() {
			var p Package
			res.Scan(&p.DateReceived, &p.DatePickedUp)
			diff := now.Sub(p.DatePickedUp)
			if diff.Hours() >= 672 {	// Check to see if the packages is at least 2 weeks old
				// Remove the pacakge from the archive
				delstmt, err := db.Prepare("DELETE FROM Picked_Up WHERE date_received=? AND date_picked_up=?")
				if err != nil {
					log.Println("Error occured while preparing statement:", err)
					return err
				}
				defer delstmt.Close()

				_, err = delstmt.Exec(p.DateReceived, p.DatePickedUp)
				if err != nil {
					log.Println("Error occured while executing statement:", err)
					return err
				}
			}
		} else {
			return ErrNoPackageFound
		}
	}

	return nil
}