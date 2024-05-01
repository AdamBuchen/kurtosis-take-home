package main

import (
	"io"
	"os"
)

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
