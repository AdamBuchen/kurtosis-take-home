package main

import (
	"fmt"
	"io"
	"os"
)

func handleFatalError(errStr string) {
	errCode := 1
	fmt.Println("error detected: " + errStr)

	os.Exit(errCode)
}

func getStringFromPath(inputPath string) (string, error) {

	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func writeStringToFile(contents string, outputPath string) error {
	err := os.WriteFile(outputPath, []byte(contents), 0666)
	if err != nil {
		return err
	}
	return nil
}
