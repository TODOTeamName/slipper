package main

import (
	"net/http"
)

func getBuilding(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", false
	}

	return sessions[cookie.Value]["building"], true
}
