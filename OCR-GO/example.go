package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
	"os/exec"
	"bytes"
	
	
)

func readLine(file *os.File) [10]string{
	var packingInfo[10] string
	var count int = 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		packingInfo[count] = scanner.Text()
		fmt.Println(packingInfo[count])
		count = count +1
    }
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
    }
	
	return packingInfo
}

func main() {
	cmd := exec.Command("tesseract", "-h")
	var out bytes.Buffer
	cmd.Stdout = &out
	err:= cmd.Run()
    file, err := os.Open("C:/Users/Quinn/Documents/School/Team soft/project/Go/out6.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	var info[10] string = readLine(file)
	fmt.Println(info[1])
	fmt.Printf("%q\n", out.String())
}