package main

type Config struct {
	// Host of the application.
	// Recommended: ":8080"
	Host *string

	// The root of where the web files are stored.
	// Recommended: "web/"
	Root *string

	// Authentication Key used for Cookie Session Storage. Should
	// be either 32 or 64 bytes long (encoded in UTF8)
	CookieKey *string

	// The https info for the website.
	// To use http instead of https (not recommended) put null here.
	Https *struct {
		// The location of the server's HTTPS cert.
		Cert *string

		// The location of the server's HTTPS key.
		Key *string
	}
}
