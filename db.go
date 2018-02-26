//http://go-database-sql.org/
package main

import "github.com/mattn/go-sqlite3"
import "log"

var db, err = nil

func Connect(){
	db, err := go-sqlite3.Open("sqlite3", "SERVER_NAME")
	if err != nil{
		log.Fatal(err)
	}
}

func Close(){
	db.Close()
}

func AddPackage(name string, building string, room string, package_type string){
	stmt, err := db.Prepare("INSERT INTO Packages(name, building, room, package_type) VALUES(?,?,?,?)")
	if err != nil{
		log.Fatal(err)
	}
	res, err := stmt.Exec(name, building, room, package_type)
	if err != nil{
		log.Fatal(err)
	}
}

