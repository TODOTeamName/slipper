package main

import (
	"./db"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var store *sessions.CookieStore

func initCookieStore() {
	store = sessions.NewCookieStore([]byte(*Settings.CookieKey))
}

func handlePackageAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	err := db.AddPackage(form.Get("name"), form.Get("building"), form.Get("room"), form.Get("type"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Database call went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}
}

func handlePackageRemove(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	err := db.RemovePackage(form.Get("number"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Database call went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}
}
