package printing

import (
	"log"
	"github.com/desertbit/fillpdf"
)

func createSlips() err{
	// Get packages to be printed

	// Add package info to the form field
	form := fillpdf.Form{
		"roomNumber1": "Room Number",
		"date1": "Date",
		"name1": "Name",
		"sortingNumber1": "Sorting Number",
		"roomNumber2": "Room Number",
		"date2": "Date",
		"name2": "Name",
		"sortingNumber2": "Sorting Number",
		"roomNumber3": "Room Number",
		"date3": "Date",
		"name3": "Name",
		"sortingNumber3": "Sorting Number",
		"roomNumber4": "Room Number",
		"date4": "Date",
		"name4": "Name",
		"sortingNumber4": "Sorting Number",
	}

	// Fill the form PDF with our values.
	err := fillpdf.Fill(form, "PackageSlipTemplateFillable.pdf", "FilledPackageSlip.pdf", true)
	if err != nil {
		log.Fatal(err)
	}
}
