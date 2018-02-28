package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var Settings Config

func main() {
	// Read Config
	configBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error while reading config: ", err)
	}

	// Parse Config JSON into struct
	err = json.Unmarshal(configBytes, &Settings)
	if err != nil {
		log.Fatal("Error while parsing config: ", err)
	}

	initCookieStore()

	// Start HTTP Server :)
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/", defHandler);
	if (Settings.Https != nil) {
		https := Settings.Https
		if (https.Cert == nil || https.Key == nil) {
			log.Fatal("Cert/Key cannot be null!");
		}

		log.Fatal(http.ListenAndServeTLS(
			Settings.Host,
			https.Cert,
			https.Key,
			nil,
		))
	} else {
		log.Fatal(http.ListenAndServe(Settings.Host, nil))
	}
}
