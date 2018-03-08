package db

import (
	"fmt"
	"log"
	"time"
)

const letters = "ABCDEFGHJKLMNPQRTUVWXY"

type Package struct {
	Number      SortingNumber
	Name        string
	Building    string
	Room        string
	PackageType string
}

type SortingNumber struct {
	Letter rune
	Number uint16
}

func getNextSortingNumber() (SortingNumber, error) {
	since := time.Since(time.Unix(0, 0))
	days := int(since.Hours()) / 24
	letter := []rune(letters)[days%len(letters)]

	res, err := db.Query(`
		SELECT MAX(CAST(SUBSTR(sorting_number, 2, 4) AS INTEGER)) FROM Packages
		WHERE SUBSTR(sorting_number, 1, 1) = ?`, letter)
	if err != nil {
		log.Println("Error while running query:", err)
		return SortingNumber{}, err
	}
	defer res.Close()

	if res.Next() {
		num := SortingNumber{}
		num.Letter = letter
		res.Scan(&num.Number)
		log.Println(num)
		num.Number++
		log.Println(num)
		return num, nil
	}

	return SortingNumber{letter, 0}, nil
}

func (s SortingNumber) String() string {
	return fmt.Sprintf("%c%04d", s.Letter, s.Number)
}