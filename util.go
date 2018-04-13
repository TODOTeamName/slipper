package main

import "net/http"

func getBuilding(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("building")
	if err != nil {
		newCookie := http.Cookie{
			Name: "building",
			Value: "Wadsworth",
		}
		http.SetCookie(w, newCookie)
		return "Wadsworth"
	}

	return cookie.Value
}
