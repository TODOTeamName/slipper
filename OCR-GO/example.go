package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func readLine(file *os.File) [10]string {
	var packingInfo [10]string
	var count int = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		packingInfo[count] = scanner.Text()
		fmt.Println(packingInfo[count])
		count = count + 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return packingInfo
}

func main() {
	cmd := exec.Command("tesseract", "cropped3.png", "out")
	cmd.Run()
	file, err := os.Open("C:/Users/Quinn/Documents/School/Team soft/project/slipper/OCR-GO/out.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var info [10]string = readLine(file)
	fmt.Println(info)
}
