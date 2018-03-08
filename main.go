package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"github.com/todoteamname/slipper/db"
)

var Settings Config

func main() {
	log.Println("Loading config...")

	ex, err := os.Executable()
	if err != nil {
		log.Println("Can't find executable location, defaulting to /srv/slipper/slipper")
		ex = "/srv/slipper/slipper"
	}

	// Read Config
	configBytes, err := ioutil.ReadFile(path.Join(path.Dir(ex), "config.json"))
	if err != nil {
		log.Fatalln("Error while reading config: ", err)
	}

	// Parse Config JSON into struct
	err = json.Unmarshal(configBytes, &Settings)
	if err != nil {
		log.Fatalln("Error while parsing config: ", err)
	}

	db.Init(path.Join(path.Dir(ex), "database.sqlite3"))

	if !path.IsAbs(*Settings.Root) {
		log.Fatalln("Root in config should be an absolute path!")
	}

	log.Println("Initializing cookie store...")
	initCookieStore()

	// Start HTTP Server :)
	http.Handle("/", http.FileServer(http.Dir(*Settings.Root)))
	http.HandleFunc("/addpackage", handlePackageAdd)
	http.HandleFunc("/removepackage", handlePackageRemove)

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
