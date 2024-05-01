package main

import (
	"flag"
	"log"
	"strings"
)

func main() {

	// Read in command-line flags, which should be our yml input and txt output paths, respectively
	flag.Parse()
	inputPath := flag.Arg(0)
	outputPath := flag.Arg(1)
	if inputPath == "" || outputPath == "" {
		log.Fatalf("two input arguments required")
	}

	// Read in Yaml string from input path
	yamlStr, fileReadErr := getStringFromPath(inputPath)
	if fileReadErr != nil {
		log.Fatal("could not open input path: " + fileReadErr.Error())
	}

	// Actually process the user job. If successful, gets back a string that can be directly inserted into
	// the file at outputPath. This is where the heavy lifting is, and there's a clear interface (YAML in, output text out)
	// that this is also where our testing can happen, at this interface boundary.
	outputLines, processingErr := ProcessUserJob(yamlStr)
	if processingErr != nil {
		log.Fatalf("could not process user job: " + processingErr.Error())
	}

	// Per instructions: An output ordering is always terminated by a newline
	var builder strings.Builder
	for _, line := range outputLines {
		builder.WriteString(line + "\n")
	}

	outputStr := builder.String()

	// Write the resulting lines of text to the file at outputPath
	saveErr := writeStringToFile(outputStr, outputPath)
	if saveErr != nil {
		log.Fatalf("could not write out file: " + saveErr.Error())
	}

	// We're done! For now, output the outputStr to STDOUT
	// TODO: Remove this

}
