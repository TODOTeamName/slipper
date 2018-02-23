//http://go-database-sql.org/
package main

import "github.com/mattn/go-sqlite3"

func addPackage( ){
	db, err := go-sqlite3.open("sqlite3", "SERVER_NAME");
	if err !+ nil{
		log.Fatal(err);
	}
	db.close();
}
