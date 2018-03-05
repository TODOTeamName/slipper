package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"path"
	"io/ioutil"
	"log"
)

var store *sessions.CookieStore

func initCookieStore() {
	store = sessions.NewCookieStore([]byte(*Settings.CookieKey))
}

// Function called when someone uses the /api/* endpoint
func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("API Call:", r.RemoteAddr, r.URL.Path)
	fmt.Fprintf(w, "<body>Hello, %s!</body>\n", r.URL.Path)
}

// Default function called when someone makes a request to the webserver
func defHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Default Handler:", r.RemoteAddr, r.URL.Path)
	
	fPath := path.Join(*Settings.Root, r.URL.Path)
	if fPath == "web" {
		fPath = "web/index.html"
	}
	bytes, err := ioutil.ReadFile(fPath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintln(w, "Error 404: File not found!")
		fmt.Fprintln(w, "Specific error: ", err)
	}
	w.Write(bytes);
	/*
	tmpl, err := template.ParseFiles(fPath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintln(w, "Error 404: File Not Found")
		fmt.Fprintln(w, "Specific error: ", err)
	}
 
	sess, err := store.Get(r, "data")
	sess.Values["user"].(string)

	//TODO: Actually replace stuff
	err = tmpl.Execute(w, struct{}{})
	*/

	if err != nil {
		log.Println("Error:", err)
	}
}
