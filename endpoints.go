package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

var store *sessions.CookieStore

func initCookieStore() {
	store = sessions.NewCookieStore([]byte(*Settings.CookieKey))
}

func handlePackageAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	r.PostForm.Get("name")
}

// Default function called when someone makes a request to the webserver
func defHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Default Handler:", r.RemoteAddr, r.URL.Path)

	if strings.Contains(r.URL.Path, "..") {
		w.WriteHeader(404)
		fmt.Fprintln(w, "Error 404: File not found!")
		fmt.Fprintln(w, "Specific error: Do not use `..` in a URL path!")
	}

	fPath := path.Join(*Settings.Root, r.URL.Path)
	if r.URL.Path == "" || r.URL.Path == "/" {
		fPath = path.Join(*Settings.Root, "index.html")
	}
	bytes, err := ioutil.ReadFile(fPath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintln(w, "Error 404: File not found!")
		fmt.Fprintln(w, "Specific error: ", err)
	}
	w.Write(bytes)

	if err != nil {
		log.Println("Error:", err)
	}
}
