// Test code dump from local testing

package printing

import (
	"bytes"
	"fmt"
	"github.com/desertbit/fillpdf"
	"github.com/todoteamname/slipper/db"
	"os/exec"
	"path"
	"strings"
	"strconv"
)

// Generate the pdf file contanining the slips to be printed
func CreateSlips(building string, root string) error {
	// Get packages to be printed
	packagesToBePrinted, err := db.GetToBePrinted(building)
	if err != nil {
		return err
	}

	// Generate slips for all the packages (4 slips per pdf)
	numPackages := len(packagesToBePrinted) // Number of packages to be printed
	numFiles := ((numPackages - 1) / 4) + 1 // Dean said this works
	pdfFiles := make([]string, numFiles)    // Slice containing all the pdf file names
	packageNum := 0

	for fileNum := 0; fileNum < numFiles; fileNum++ {
		// Generate the pdf file name
		fileName := fmt.Sprintf("packageSlip%03d.pdf", fileNum)

		form := fillpdf.Form{}

		for i := 1; i <= 4; i++ {
			iStr := strconv.Itoa(i)

			if packageNum >= numPackages {
				form["roomNumber" + iStr] = ""
				form["date" + iStr] = ""
				form["name" + iStr] = ""
				form["sortingNumber" + iStr] = ""
				form["carrier" + iStr] = ""
				form["packageType" + iStr] = ""

				packageNum++
				continue
			}

			form["roomNumber" + iStr] = packagesToBePrinted[packageNum].Room
			form["date" + iStr] = packagesToBePrinted[packageNum].DateReceived.Format("Mon Jan _2 3:04PM")
			form["name" + iStr] = packagesToBePrinted[packageNum].Name
			form["sortingNumber" + iStr] = packagesToBePrinted[packageNum].Number.String()
			form["carrier" + iStr] = packagesToBePrinted[packageNum].Carrier
			form["packageType" + iStr] = packagesToBePrinted[packageNum].PackageType

			packageNum++
		}


		// Fill the form PDF with our values.
		err = fillpdf.Fill(form, path.Join(root, "../printing/PackageSlipTemplate.pdf"), path.Join(root, fileName), true)
		if err != nil {
			return err
		}
		pdfFiles[fileNum] = path.Join(root, fileName)
	}

	// Collate all pdf files togethers
	args := make([]string, numFiles+3)
	var argNum int
	for argNum = 0; argNum < numFiles; argNum++ {
		args[argNum] = pdfFiles[argNum]
	}
	args[argNum] = "cat"
	argNum++
	args[argNum] = "output"
	argNum++
	args[argNum] = path.Join(root, "PackageSlips.pdf")

	var stderr bytes.Buffer
	cmd := exec.Command("pdftk", args...)
	cmd.Stderr = &stderr
	cmd.Dir = root

	// Start the command and wait for it to exit.
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}

	// Mark the packages as printed in the db
	err = db.MarkPrinted(building)
	if err != nil {
		return err
	}

	return nil
}
