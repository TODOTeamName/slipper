package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"path"
)

var store *sessions.CookieStore

func initCookieStore() {
	store = sessions.NewCookieStore(Config.CookieKey)
}

// Function called when someone uses the /api/* endpoint
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<body>Hello, %s!</body>\n", r.URL.Path)
}

// Default function called when someone makes a request to the webserver
func defHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := path.Join(Config.Root, r.URL.Path)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintln(w, "Error 404: File Not Found")
		fmt.Fprintln(w, "Specific error: ", err)
	}

	sess, err := store.Get(r, "data")
	username := sess.Values["user"].(string)

	err = tmpl.Execute(w, username)
}
