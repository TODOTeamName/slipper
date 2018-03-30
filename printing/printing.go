// Test code dump from local testing

package printing

import (
	"github.com/desertbit/fillpdf"
	"github.com/todoteamname/slipper/db"
	"path"
	"fmt"
	"io/ioutil"
)

var PackageSlipsPdf string

// Generate the pdf file contanining the slips to be printed
func CreateSlips(building string, root string) error {
	// Get packages to be printed
	packagesToBePrinted, err := db.GetToBePrinted(building)
	if err != nil {
		return err
	}

	/*

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

	*/

	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "packageSlips-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}

	// Remove the temporary directory on defer again.
	defer func() {
		errD := os.RemoveAll(tmpDir)
		// Log the error only.
		if errD != nil {
			log.Printf("fillpdf: failed to remove temporary directory '%s' again: %v", tmpDir, errD)
		}
	}()

	// Generate slips for all the packages (4 slips per pdf)
	numPackages := 	len(packagesToBePrinted)		// Number of packages to be printed
	numFiles 	:= 	((numPackages - 1) / 4) + 1		// Dean said this works
	pdfFiles 	:= 	make([]string, numFiles)		// Slice containing all the pdf file names
	packageNum 	:= 	0								// Counter to track which package is being processed

	for fileNum := 0; fileNum < numFiles; fileNum++{
		// Generate the pdf file name
		fileName := fmt.Sprintf("packageSlip%d.pdf", fileNum)

		// Popluate pacakge information into variables
		if packageNum < numPackages{
			roomNumber1 := packagesToBePrinted[packageNum].Room
			date1 := packagesToBePrinted[packageNum].DateReceived
			name1 := packagesToBePrinted[packageNum].Name
			sortingNumber1 := packagesToBePrinted[packageNum].Number.String()
			carrier1 := packagesToBePrinted[packageNum].Carrier
			packageType1 := packagesToBePrinted[packageNum].PackageType
		}else{
			roomNumber1 := ""
			date1 := ""
			name1 := ""
			sortingNumber1 := ""
			carrier1 := ""
			packageType1 := ""
		}
		packageNum++

		if packageNum < numPackages{
			roomNumber2 := packagesToBePrinted[packageNum].Room
			date2 := packagesToBePrinted[packageNum].DateReceived
			name2 := packagesToBePrinted[packageNum].Name
			sortingNumber2 := packagesToBePrinted[packageNum].Number.String()
			carrier2 := packagesToBePrinted[packageNum].Carrier
			packageType2 := packagesToBePrinted[packageNum].PackageType
		}else{
			roomNumber2 := ""
			date2 := ""
			name2 := ""
			sortingNumber2 := ""
			carrier2 := ""
			packageType2 := ""
		}
		packageNum++

		if packageNum < numPackages{
			roomNumber3 := packagesToBePrinted[packageNum].Room
			date3 := packagesToBePrinted[packageNum].DateReceived
			name3 := packagesToBePrinted[packageNum].Name
			sortingNumber3 := packagesToBePrinted[packageNum].Number.String()
			carrier3 := packagesToBePrinted[packageNum].Carrier
			packageType3 := packagesToBePrinted[packageNum].PackageType
		}else{
			roomNumber3 := ""
			date3 := ""
			name3 := ""
			sortingNumber3 := ""
			carrier3 := ""
			packageType3 := ""
		}
		packageNum++

		if packageNum < numPackages{
			roomNumber4 := packagesToBePrinted[packageNum].Room
			date4 := packagesToBePrinted[packageNum].DateReceived
			name4 := packagesToBePrinted[packageNum].Name
			sortingNumber4 := packagesToBePrinted[packageNum].Number.String()
			carrier4 := packagesToBePrinted[packageNum].Carrier
			packageType4 := packagesToBePrinted[packageNum].PackageType
		}else{
			roomNumber4 := ""
			date4 := ""
			name4 := ""
			sortingNumber4 := ""
			carrier4 := ""
			packageType4 := ""
		}
		packageNum++

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
		err = fillpdf.Fill(form, path.Join(root, "../printing/PackageSlipTemplate.pdf"), path.Join(root, tmpDir, fileName, true)
		if err != nil {
			return err
		}
		pdfFiles[fileNum] = path.Join(root, tmpDir, fileName, true)
	}

	// Collate all pdf files togethers
	args := make([]string, numFiles+3)
	for argNum := 0; argNum < numFiles; argNum++{
		args[argNum] = pdfFiles[argNum]
	}
	args[argNum++] = "cat"
	args[argNum++] = "output"
	args[argNum++] = path.Join(root, "PackageSlips", true)
	err = fillpdf.runCommandInPath(tmp, "pdftk", args...)
	if err != nil {
		return fmt.Errorf("pdftk error: %v", err)
	}

	// Mark the packages as printed in the db
	return nil
}
