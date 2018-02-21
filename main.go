package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var Settings Config

func main() {
	// Read Config
	configBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error while reading config: ", err)
	}

	// Parse Config JSON into struct
	err = json.Unmarshal(configBytes, &Config)
	if err != nil {
		log.Fatal("Error while parsing config: ", err)
	}

	// Start HTTP Server :)
	http.HandleFunc("/", handler)
	if (Config.Https != nil) {
		https := Config.Https
		if (https.Cert == nil || https.Key == nil) {
			log.Fatal("Cert/Key cannot be null!");
		}

		log.Fatal(http.ListenAndServeTLS(
			Config.Host,
			https.Cert,
			https.Key,
			nil
		))
	} else {
		log.Fatal(http.ListenAndServe(Config.Host, nil))
	}
}

// Handler function to handle HTTP Requests
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<body>lol</body>")
}
