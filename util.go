package main

import (
	"net/http"
)

// Gets the building information from the session token.
//
// Returns the building that a request is logged into,
// and a boolean indicating if it was successful.
func getBuilding(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", false
	}
	if _, ok := sessions[cookie.Value]; !ok {
		return "", false
	}

	return sessions[cookie.Value]["building"], true
}
