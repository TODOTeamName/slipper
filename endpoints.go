package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/todoteamname/slipper/db"
	"github.com/todoteamname/slipper/ocr"
	"github.com/todoteamname/slipper/printing"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"text/template"
	"math/rand"
)

var sessions = make(map[string]map[string]string)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	building := r.FormValue("building")
	pass := r.FormValue("password")

	hash, err := db.GetPassword(building)
	if err != nil {
		fmt.Fprintf(w, "<body>Database error: %s</body>", err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		fmt.Fprintf(w, "<body>Invalid login, try again... (redirecting)</body>")
		fmt.Fprintf(w, "<script>setTimeout(function() { window.location='/' }, 3000)</script>")
		return
	}

	// select a session id (done secure-ish-ly)
	sessid := strconv.Itoa(rand.Int())
	for _, ok := sessions[sessid]; ok; _, ok = sessions[sessid] {
		sessid = strconv.Itoa(rand.Int())
	}

	c := http.Cookie{Name:"session", Value:sessid}
	http.SetCookie(w, &c)

	sessions[sessid] = map[string]string {"building":building}

	fmt.Fprintf(w, "<body>Login successful! (redirecting)</body>")
	fmt.Fprintf(w, "<script>setTimeout(function() { window.location='/pages/main.html' }, 3000)</script>")
}

func handlePackageAdd(w http.ResponseWriter, r *http.Request) {

	building, ok := getBuilding(r)
	if !ok {
		w.WriteHeader(401)
		fmt.Fprintln(w, "Error 401: Forbidden. Please log in.")
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	num, err := db.AddPackage(
		r.FormValue("name"),
		building,
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

	building, ok := getBuilding(r)
	if !ok {
		w.WriteHeader(401)
		fmt.Fprintln(w, "Error 401: Forbidden. Please log in.")
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	sigB64 := r.FormValue("sig")
	sig, err := base64.RawStdEncoding.DecodeString(sigB64)

	err = db.Archive(r.FormValue("number"), building, sig)
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

	building, ok := getBuilding(r)
	if !ok {
		w.WriteHeader(401)
		fmt.Fprintln(w, "Error 401: Forbidden. Please log in.")
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	pack, err := db.GetPackage(r.FormValue("number"), building)
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

func handlePackageUpdate(w http.ResponseWriter, r *http.Request) {
	building, ok := getBuilding(r)
	if !ok {
		w.WriteHeader(401)
		fmt.Fprintln(w, "Error 401: Forbidden. Please log in.")
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	isPrinted, _ := strconv.Atoi(r.FormValue("isprinted"))

	err := db.UpdatePackage(
		r.FormValue("sortingnumber"),
		r.FormValue("name"),
		building,
		r.FormValue("room"),
		r.FormValue("carrier"),
		r.FormValue("type"),
		isPrinted,
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

	building, ok := getBuilding(r)
	if !ok {
		w.WriteHeader(401)
		fmt.Fprintln(w, "Error 401: Forbidden. Please log in.")
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	err := printing.CreateSlips(building, *Settings.Root)
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
	cmd := exec.Command("rm", "*.pdf")
	cmd.Stderr = &stderr
	cmd.Dir = *Settings.Root
	err = cmd.Run()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Something went wrong in the removal.")
		fmt.Fprintln(w, "Precise error:", err)
		return
	}
}

func handleOcr(w http.ResponseWriter, r *http.Request) {
	_, ok := getBuilding(r)
	if !ok {
		w.WriteHeader(401)
		fmt.Fprintln(w, "Error 401: Forbidden. Please log in.")
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}
	r.ParseMultipartForm(500000)
	fmt.Println(r.MultipartForm.Value)
	b64 := r.FormValue("baseimage")
	byt, err := base64.RawStdEncoding.DecodeString(b64)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. Decoding went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}
	output, err := ocr.ReadFile(byt)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Error 400: Bad Request. OCR went wrong.")
		fmt.Fprintln(w, "Precise error:", err)
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	fmt.Fprintf(w, "OCR Output: %s", output)
}

func handleCheckArchive(w http.ResponseWriter, r *http.Request) {
	building, ok := getBuilding(r)
	if !ok {
		w.WriteHeader(401)
		fmt.Fprintln(w, "Error 401: Forbidden. Please log in.")
		fmt.Fprintln(w, "Click <a href=\"/\">here</a> to go to the home page")
		return
	}

	archivePackages := db.checkArchive(r.FormValue("name"), building)
	
	
	return
}
