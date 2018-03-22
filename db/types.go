package db

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"errors"
)

const letters = "ABCDEFGHJKLMNPQRTUVWXY"

type Package struct {
	Number      *SortingNumber
	Name        string
	Building    string
	Room        string
	PackageType string
}

type SortingNumber struct {
	Letter rune
	Number uint16
}

func getNextSortingNumber() (*SortingNumber, error) {
	since := time.Since(time.Unix(0, 0))
	days := int(since.Hours()) / 24
	letter := []rune(letters)[days%len(letters)]

	res, err := db.Query(`
		SELECT MAX(CAST(SUBSTR(sorting_number, 2, 4) AS INTEGER)) FROM Packages
		WHERE SUBSTR(sorting_number, 1, 1) = ?`, string(letter))
	if err != nil {
		log.Println("Error while running query:", err)
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		num := new(SortingNumber)
		num.Letter = letter
		res.Scan(&num.Number)
		num.Number++
		return num, nil
	}

	num := SortingNumber{letter, 0}
	return &num, nil
}

func (s *SortingNumber) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("uh oh")
	}

	fmt.Println("Scan called with", src)
	s.Letter = []rune(str)[0]
	conv, _ := strconv.Atoi(str[1:])
	s.Number = uint16(conv)
	fmt.Println("Scan call mutates to", s)
	return nil
}

func (s *SortingNumber) String() string {
	return fmt.Sprintf("%c%04d", s.Letter, s.Number)
}
