package ocr

import (
	"os/exec"
)

func ReadFile(input []byte) (string, error) {
	cmd := exec.Command("tesseract", "stdin", "stdout")
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	_, err = stdin.Write(input)
	if err != nil {
		return "", err
	}

	err = stdin.Close()
	if err != nil {
		return "", err
	}

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return bytes, nil
}
