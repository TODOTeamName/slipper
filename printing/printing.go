// Test code dump from local testing

package printing

import (
	"github.com/desertbit/fillpdf"
	"github.com/todoteamname/slipper/db"
	"path"
)

var PackageSlipsPdf string

// Generate the pdf file contanining the slips to be printed
func CreateSlips(building string, root string) error {
	// Get packages to be printed
	packagesToBePrinted, err := db.GetToBePrinted(building)
	if err != nil {
		return err
	}

	// Set up the form values for the first 4 packages
	roomNumber1 := packagesToBePrinted[0].Room
	date1 := packagesToBePrinted[0].DateReceived
	name1 := packagesToBePrinted[0].Name
	sortingNumber1 := packagesToBePrinted[0].Number.String()
	carrier1 := packagesToBePrinted[0].Carrier
	packageType1 := packagesToBePrinted[0].PackageType

	roomNumber2 := packagesToBePrinted[1].Room
	date2 := packagesToBePrinted[1].DateReceived
	name2 := packagesToBePrinted[1].Name
	sortingNumber2 := packagesToBePrinted[1].Number.String()
	carrier2 := packagesToBePrinted[1].Carrier
	packageType2 := packagesToBePrinted[1].PackageType

	roomNumber3 := packagesToBePrinted[2].Room
	date3 := packagesToBePrinted[2].DateReceived
	name3 := packagesToBePrinted[2].Name
	sortingNumber3 := packagesToBePrinted[2].Number.String()
	carrier3 := packagesToBePrinted[2].Carrier
	packageType3 := packagesToBePrinted[2].PackageType

	roomNumber4 := packagesToBePrinted[3].Room
	date4 := packagesToBePrinted[3].DateReceived
	name4 := packagesToBePrinted[3].Name
	sortingNumber4 := packagesToBePrinted[3].Number.String()
	carrier4 := packagesToBePrinted[3].Carrier
	packageType4 := packagesToBePrinted[3].PackageType

	// Add package info to the form field
	form := fillpdf.Form{
		"roomNumber1":    roomNumber1,
		"date1":          date1,
		"name1":          name1,
		"sortingNumber1": sortingNumber1,
		"carrier1":       carrier1,
		"packageType1":   packageType1,
		"roomNumber2":    roomNumber2,
		"date2":          date2,
		"name2":          name2,
		"sortingNumber2": sortingNumber2,
		"carrier2":       carrier2,
		"packageType2":   packageType2,
		"roomNumber3":    roomNumber3,
		"date3":          date3,
		"name3":          name3,
		"sortingNumber3": sortingNumber3,
		"carrier3":       carrier3,
		"packageType3":   packageType3,
		"roomNumber4":    roomNumber4,
		"date4":          date4,
		"name4":          name4,
		"sortingNumber4": sortingNumber4,
		"carrier4":       carrier4,
		"packageType4":   packageType4,
	}

	// Fill the form PDF with our values.
	err = fillpdf.Fill(form, path.Join(root, "../printing/PackageSlipTemplate.pdf"), path.Join(root, "FilledPackageSlip.pdf"), true)
	if err != nil {
		return err
	}

	// Mark the packages as printed in the db
	return nil
}

func printSlips() {

}
