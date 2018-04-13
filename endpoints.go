package main

import (
	"bytes"
	"fmt"
	"github.com/todoteamname/slipper/db"
	"github.com/todoteamname/slipper/printing"
	"net/http"
	"path"
	"text/template"
	"os"
	"os/exec"
	"io"
	"io/ioutil"
	"github.com/todoteamname/slipper/ocr"
	"strconv"
)

func handlePackageAdd(w http.ResponseWriter, r *http.Request) {

	num, err := db.AddPackage(
		r.FormValue("name"),
		r.FormValue("building"),
		r.FormValue("room"),
		r.FormValue("carrier"),
		r.FormValue("type"),
	)
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

	err := db.Archive(r.FormValue("number"), r.FormValue("signature"))
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

	pack, err := db.GetPackage(r.FormValue("number"))
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

func handlePackageUpdate(w http.ResponseWriter, r *http.Request){
	isPrinted, _ := strconv.Atoi(r.FormValue("isprinted"))

	err := db.UpdatePackage(
		r.FormValue("sortingnumber"),
		r.FormValue("name"),
		r.FormValue("building"),
		r.FormValue("room"),
		r.FormValue("carrier"),
		r.FormValue("type"),
		isPrinted
	)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Database call went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}
}

func handleCreateSlips(w http.ResponseWriter, r *http.Request) {

	err := printing.CreateSlips(r.FormValue("building"), *Settings.Root)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Assembling PDF went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	

	f, err := os.Open(path.Join(*Settings.Root, "PackageSlips.pdf"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	// Set header
	w.Header().Add("Content-Type", "application/pdf")

	// Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}

	// Remove the package slip files
	var stderr bytes.Buffer
	args := make([]string, 1)
	args[0] = "*.pdf"
	cmd := exec.Command("rm", args...)
	cmd.Stderr = &stderr
	cmd.Dir = *Settings.Root
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(w, "Error 400: Something went wrong in the removal.")
		fmt.Fprintln(w, "Precise error:", err)
	}
}

func handleOcr(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(500000)
	fmt.Println(r.MultipartForm.Value)
	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Form Parsing went wrong")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}
	defer file.Close()

	fb, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Reading File went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}
	output, err := ocr.ReadFile(fb)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. OCR went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	fmt.Fprintf(w, "OCR Output: %s", output)
}
