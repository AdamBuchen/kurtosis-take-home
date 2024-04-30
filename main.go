package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"gopkg.in/yaml.v3"
)

func main() {

	// Read in command-line flags, which should be our yml input and txt output paths, respectively
	flag.Parse()
	inputPath := flag.Arg(0)
	outputPath := flag.Arg(1)

	fmt.Println(inputPath)
	fmt.Println(outputPath)

	// 1. Read in the YAML file as a string. This will make testing much easier and separate file io from content
	// Get input path and output path flags

	// TEST:

	testYamlStr := `
- step: "create user 1"
  dependencies: ["prepare database"]
  precedence: 100
- step: "create user 2"
  dependencies: ["prepare database"]
  precedence: 50
- step: "prepare database"
  dependencies: []
  precedence: 10
- step: "create user 3"
  dependencies: ["create user 4"]
  precedence: 10
- step: "create user 4"
  dependencies: ["create user 2"]
  precedence: 100
`
	// 2. Feed the string into Go YAML parser to get back an InputJob
	inputJob := InputJob{}
	yamlMarshalErr := yaml.Unmarshal([]byte(testYamlStr), &inputJob.Steps)
	if yamlMarshalErr != nil {
		log.Fatalf("invalid yaml: %s", yamlMarshalErr)
	}

	// 3. Take the inputJob.Steps and get stepsByIdMap (also do validation here)
	stepsByIdMap, stepsErr := getStepsByIdMap(inputJob.Steps)
	if stepsErr != nil {
		log.Fatalf("could not get steps: %s", stepsErr)
	}

	// 4. Take the stepsByIdMap and get a depsByParent map
	depsByParent, depsByParentErr := getDepsByParent(stepsByIdMap)
	if depsByParentErr != nil {
		log.Fatalf("could not get deps by parent: %s", depsByParentErr)
	}

	// 5. Loop through stepsByIdMap. Recursively run processChildItem on each of them. By the end, every item should have step group number, and precedence, and name
	currentGroupNumber := 1
	for _, step := range stepsByIdMap {
		processingErr := processChildItem(step, depsByParent, currentGroupNumber)
		if processingErr != nil {
			log.Fatalf("could not process items: %s", processingErr)
		}
	}

	// 6. Copy map to array, verify that all dependencies are clear
	sortingArr := make([]*JobStep, len(stepsByIdMap))
	i := 0

	for _, step := range stepsByIdMap {
		// Deps should be sorted out by now. If not, input was broken.
		if !step.AllParentDepsClear() {
			log.Fatalf("dependency issue detected")
		}

		sortingArr[i] = step
		i++
	}

	// 7. Custom sort: StepGroupNumber asc, Precedence desc, StepId asc
	sort.Slice(sortingArr, func(i, j int) bool {
		if sortingArr[i].StepGroupNumber != sortingArr[j].StepGroupNumber {
			return sortingArr[i].StepGroupNumber < sortingArr[j].StepGroupNumber
		}
		if sortingArr[i].Precedence != sortingArr[j].Precedence {
			return sortingArr[i].Precedence > sortingArr[j].Precedence
		}
		return sortingArr[i].StepId < sortingArr[j].StepId
	})

	// 8. Output the lines to a file and return status code 0

	for _, step := range sortingArr {
		fmt.Println(step.StepId)
		fmt.Println(step.StepGroupNumber)
	}

}
