package main

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

// This function is the primary interface boundary for the project. It takes
// in a string of Yaml and returns an array of strings ready to be put in the user's
// out file. Everything that's testable can be tested at this boundary.
// The logic outside of the code is primarily os / io utilities (files, flags, sys codes)
func ProcessUserJob(yamlStr string) ([]string, error) {

	output := make([]string, 0)

	// 1. Feed the string into Go YAML parser to get back an InputJob
	inputJob := InputJob{}
	yamlMarshalErr := yaml.Unmarshal([]byte(yamlStr), &inputJob.Steps)
	if yamlMarshalErr != nil {
		return output, fmt.Errorf("invalid yaml: %s", yamlMarshalErr)
	}

	// 2. Take the inputJob.Steps and get stepsByIdMap (also do validation here)
	stepsByIdMap, stepsErr := getStepsByIdMap(inputJob.Steps)
	if stepsErr != nil {
		return output, fmt.Errorf("could not get steps: %s", stepsErr)
	}

	if len(stepsByIdMap) == 0 {
		return output, fmt.Errorf("no steps were provided by user")
	}

	// 3. Take the stepsByIdMap and get a depsByParent map
	depsByParent, depsByParentErr := getDepsByParent(stepsByIdMap)
	if depsByParentErr != nil {
		return output, fmt.Errorf("could not get deps by parent: %s", depsByParentErr)
	}

	// 4. Loop through stepsByIdMap. Recursively run processChildItem on each of them. By the end, every item should have step cycle number
	currentCycleNumber := 1
	for _, step := range stepsByIdMap {
		processingErr := processChildItem(step, depsByParent, currentCycleNumber)
		if processingErr != nil {
			return output, fmt.Errorf("could not process items: %s", processingErr)
		}
	}

	// 5. Copy map to array, verify that all dependencies are clear
	sortingArr := make([]*JobStep, len(stepsByIdMap))
	i := 0

	for _, step := range stepsByIdMap {
		// Deps should be sorted out by now. If not, input was broken.
		if step.StepCycleNumber == 0 || !step.AllParentDepsClear() {
			return output, fmt.Errorf("dependency issue detected")
		}

		sortingArr[i] = step
		i++
	}

	// 6. Custom sort: StepCycleNumber asc, Precedence desc, StepId asc
	sort.Slice(sortingArr, func(i, j int) bool {
		if sortingArr[i].StepCycleNumber != sortingArr[j].StepCycleNumber {
			return sortingArr[i].StepCycleNumber < sortingArr[j].StepCycleNumber
		}
		if sortingArr[i].Precedence != sortingArr[j].Precedence {
			return sortingArr[i].Precedence > sortingArr[j].Precedence
		}
		return sortingArr[i].StepId < sortingArr[j].StepId
	})

	// 7. Loop through our array once more, build up a string.
	// Per instructions: "An output ordering is always terminated by a newline"
	for _, step := range sortingArr {
		output = append(output, step.StepId)
	}

	return output, nil
}

// Return a map which is keyed by the id of the step, given an array of user inputs
func getStepsByIdMap(inputSteps []*InputJobStep) (map[string]*JobStep, error) {

	output := make(map[string]*JobStep)
	for _, inputStep := range inputSteps {
		validationErr := inputStep.ValidateInputStep()
		if validationErr != nil {
			return output, fmt.Errorf("validation error received: " + validationErr.Error())
		}

		jobStep := inputStep.GetJobStep()

		// TODO: Check if key already exists: if so, there's a dupe, which should return error
		output[jobStep.StepId] = jobStep
	}

	// Now we loop over our dependecyIds (array of strings), for each step, and assign a pointer
	// to that actual parent node in our DepsToClear[parentStepId] map. If a node cannot be found,
	// that means that the input dependency string does not match any actual steps.
	for stepId, step := range output {
		for _, parentStepId := range step.DependencyIds {
			if parent, parentOk := output[parentStepId]; parentOk {
				output[stepId].DepsToClear[parentStepId] = parent
			} else {
				return output, fmt.Errorf("invalid dependency specified: %s", parentStepId)
			}
		}
	}

	return output, nil
}

// Returns a map keyed by parentStepId, each containing an array of direct children steps
// This is useful later when needing to do reverse lookups. When a step's parents'
// dependencies (if any) resolve, then the step is clear to get a cycle number.
// We want to know if that step had any child dependencies so we can tell them
// one of their parent dependencies just cleared. To do that lookup, we need
// this data structure.
func getDepsByParent(stepsByIdMap map[string]*JobStep) (map[string][]*JobStep, error) {

	output := make(map[string][]*JobStep)
	for stepId := range stepsByIdMap {
		output[stepId] = make([]*JobStep, 0)
	}

	for _, step := range stepsByIdMap {
		for parentId := range step.DepsToClear {
			output[parentId] = append(output[parentId], step)
		}
	}

	return output, nil
}

// A recursive function that checks whether in the current cycle we can yet
// process this dependency. Results in Cycle Number being set and AllDepsClear
// being set to true for any children that can be processed.
func processChildItem(step *JobStep, depsByParent map[string][]*JobStep, currentCycleNumber int) error {

	if step.StepCycleNumber > 0 && currentCycleNumber > step.StepCycleNumber {
		return fmt.Errorf("cyclical relationship detected: %s", step.StepId)
	}

	// This one was already processed, let's not get stuck in a loop
	if step.AllDepsClear || step.StepCycleNumber > 0 {
		return nil
	}

	if step.AllParentDepsClear() {
		step.AllDepsClear = true
		step.StepCycleNumber = currentCycleNumber
		if _, ok := depsByParent[step.StepId]; ok {
			nextCycleNumber := currentCycleNumber + 1
			for _, childStep := range depsByParent[step.StepId] {
				err := processChildItem(childStep, depsByParent, nextCycleNumber)
				if err != nil {
					return err // Just pass it up the chain
				}
			}
		}
	}

	return nil
}
