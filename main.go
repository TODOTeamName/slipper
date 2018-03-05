package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var Settings Config

func main() {
	log.Println("Loading config...")
	
	// Read Config
	configBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("Error while reading config: ", err)
	}

	// Parse Config JSON into struct
	err = json.Unmarshal(configBytes, &Settings)
	if err != nil {
		log.Fatalln("Error while parsing config: ", err)
	}

	log.Println("Initializing cookie store...")
	initCookieStore()

	// Start HTTP Server :)
	http.HandleFunc("/api/", apiHandler)

	http.HandleFunc("/", defHandler)

	log.Println("Starting server...")
	if Settings.Https != nil {
		https := Settings.Https
		if https.Cert == nil || https.Key == nil {
			log.Fatal("Cert/Key cannot be null!")
		}

		log.Fatal(http.ListenAndServeTLS(
			*Settings.Host,
			*https.Cert,
			*https.Key,
			nil,
		))
	} else {
		log.Fatal(http.ListenAndServe(*Settings.Host, nil))
	}
}
