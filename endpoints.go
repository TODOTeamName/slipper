package main

import (
	"fmt"
	"github.com/todoteamname/slipper/db"
	"net/http"
	"path"
)

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

	fmt.Fprintf(w,
		"<script>history.replaceState(%q, %q, %q);</script>",
		"asdf",
		"Slipper|Add Package",
		"/pages/form_add.html",
	)
	fmt.Fprintf(w, "<script>alert(%q)</script>", fmt.Sprintf("The package number is %s", num))
	http.ServeFile(w, r, path.Join(*Settings.Root, "pages/form_add.html"))
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

	fmt.Fprintf(w,
		"<script>history.replaceState(%q, %q, %q);</script>",
		"asdf",
		"Slipper|Remove Package",
		"/pages/form_remove.html",
	)
	http.ServeFile(w, r, path.Join(*Settings.Root, "pages/form_remove.html"))
}

func handlePackageGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	pack, err := db.GetPackage(form.Get("number"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Database call went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	fmt.Fprintf(w,
		"<script>history.replaceState(%q, %q, %q);</script>",
		"asdf",
		"Slipper|Update Package",
		"/pages/form_update.html",
	)

	fmt.Fprintf(w, "Package info: %#v<br><br>", pack)
	fmt.Fprintln(w, "This is a temporary page")
}