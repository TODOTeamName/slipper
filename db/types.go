package db

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"unicode"
)

// All letters used for the sorting number's letters
const letters = "ABCDEFGHJKLMNPQRTUWXY"

// Defines a package
type Package struct {
	Number       SortingNumber
	DateReceived time.Time
	Name         string
	Building     string
	Room         string
	Carrier      string
	PackageType  string
	Printed      bool
	DatePickedUp time.Time
	Signature	  string
}

// Represents a sorting number.
type SortingNumber struct {
	Letter rune
	Number uint16
}

// Gets the next sorting number. Note: does a database call.
func getNextSortingNumber(building string) (SortingNumber, error) {
	since := time.Since(time.Unix(0, 0))
	days := int(since.Hours()) / 24
	letter := []rune(letters)[days%len(letters)]

	res, err := db.Query(`
		SELECT MAX(CAST(SUBSTR(sorting_number, 2, 4) AS INTEGER)) FROM Packages
		WHERE SUBSTR(sorting_number, 1, 1) = ? AND building = ?`, string(letter), building)
	if err != nil {
		log.Println("Error while running query:", err)
		return SortingNumber{}, err
	}
	defer res.Close()

	if res.Next() {
		var num SortingNumber
		num.Letter = letter
		res.Scan(&num.Number)
		num.Number++
		return num, nil
	}

	return SortingNumber{letter, 0}, nil
}

// Turns a string into a sorting number
func Atosn(src string) SortingNumber {
	var val SortingNumber
	val.Letter = unicode.ToUpper([]rune(src)[0])
	conv, _ := strconv.Atoi(src[1:])
	val.Number = uint16(conv)
	if val.Letter == 'V' {
		val.Letter = 'U'
	}
	return val
}

// Turns a sorting number to a string
func (s SortingNumber) String() string {
	return fmt.Sprintf("%c%04d", s.Letter, s.Number)
}
