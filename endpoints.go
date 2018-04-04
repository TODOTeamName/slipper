package main

import (
	"fmt"
	"github.com/todoteamname/slipper/db"
	"github.com/todoteamname/slipper/printing"
	"net/http"
	"path"
	"text/template"
	/*"os"
	"io"*/
)

func handlePackageAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	num, err := db.AddPackage(form.Get("name"), form.Get("building"), form.Get("room"), form.Get("carrier"), form.Get("type"))
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

	err := db.Archive(form.Get("number"), form.Get("signature"))
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

	t, err := template.ParseFiles(path.Join(*Settings.Root, "pages/form_update.html"))
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

	t.Execute(w, pack)
}

func handleCreateSlips(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	err := printing.CreateSlips(form.Get("building"), *Settings.Root)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Assembling PDF went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	/*

	f, err := os.Open(path.Join(*Settings.Root, "PackageSlips.pdf"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//Set header
	w.Header().Add("Content-Type", "application/pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}

	*/

}
