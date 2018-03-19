package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/todoteamname/slipper/db"
	"net/http"
)

var store *sessions.CookieStore

func initCookieStore() {
	store = sessions.NewCookieStore([]byte(*Settings.CookieKey))
}

func handlePackageAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	num, err := db.AddPackage(form.Get("name"), form.Get("building"), form.Get("room"), form.Get("type"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Database call went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return

	}

	http.ServeFile(w, r, "/pages/form_add.html")
	fmt.Fprintf(w,
		"<script>history.replaceState(%q, %q, %q);</script>",
		"asdf",
		"Slipper|Add Package",
		"/pages/form_add.html",
	)
	fmt.Fprintf(w, "<script>alert(%q)</script>", fmt.Sprintf("The package number is %d", num))
}

func handlePackageRemove(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	err := db.Archive(form.Get("number"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Database call went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	http.ServeFile(w, r, "/pages/form_remove.html")
	fmt.Fprintf(w,
		"<script>history.replaceState(%q, %q, %q);</script>",
		"asdf",
		"Slipper|Remove Package",
		"/pages/form_remove.html",
	)
	fmt.Fprintln(w, "Lol it worked")
}
