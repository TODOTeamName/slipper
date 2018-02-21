package main

type Config struct {
	// Host of the application.
	// Recommended: ":8080"
	Host string

	// The root of where the web files are stored.
	// Recommended: "templates/"
	Root string

	// The https info for the website.
	// To use http instead of https (not recommended) put null here.
	Https struct {
		// The location of the server's HTTPS cert.
		Cert string

		// The location of the server's HTTPS key.
		Key string
	}
}
