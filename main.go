package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Config struct {
	Host string
}

func main() {
	// Read Config
	configBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error while reading config: ", err)
	}

	// Parse Config JSON into struct
	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal("Error while parsing config: ", err)
	}

	// Start HTTP Server :)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(config.Host, nil))
}

// Handler function to handle HTTP Requests
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<body>lol</body>")
}
