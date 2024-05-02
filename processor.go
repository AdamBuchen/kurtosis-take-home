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

	// 2. Take the inputJob.Steps and get stepsByIdSlice (also do validation here)
	stepsByIdSlice, stepsErr := getStepsByIdSlice(inputJob.Steps)
	if stepsErr != nil {
		return output, fmt.Errorf("could not get steps: %s", stepsErr)
	}

	if len(stepsByIdSlice) == 0 {
		return output, fmt.Errorf("no steps were provided by user")
	}

	// Cycle through our tree until we get nothing else
	for {
		nextStep := getNextAvailableStep(stepsByIdSlice)
		//Actually update my nextStep node

		if nextStep == nil {
			if len(output) != len(stepsByIdSlice) { // Possible circular dependency
				return output, fmt.Errorf("possible circular dependency detected")
			}
			break
		}

		output = append(output, nextStep.StepId)
	}

	return output, nil
}

// Return a map which is keyed by the id of the step, given an array of user inputs
func getStepsByIdSlice(inputSteps []*InputJobStep) ([]*JobStep, error) {

	output := make([]*JobStep, 0)
	stepsByIdMap := make(map[string]*JobStep)
	for _, inputStep := range inputSteps {
		validationErr := inputStep.ValidateInputStep()
		if validationErr != nil {
			return output, fmt.Errorf("validation error received: " + validationErr.Error())
		}

		jobStep := inputStep.GetJobStep()

		// Check if key already exists: if so, there's a dupe, which should return error
		if _, isDuplicateStep := stepsByIdMap[jobStep.StepId]; isDuplicateStep {
			return output, fmt.Errorf("duplicate key detected: %s", jobStep.StepId)
		}

		stepsByIdMap[jobStep.StepId] = jobStep
		output = append(output, jobStep)
	}

	// Now we loop over our dependecyIds (array of strings), for each step, and assign a pointer
	// to that actual parent node in our DepsToClear[parentStepId] map. If a node cannot be found,
	// that means that the input dependency string does not match any actual steps.
	for i, step := range output {
		for _, parentStepId := range step.DependencyIds {
			if parent, parentOk := stepsByIdMap[parentStepId]; parentOk {
				stepId := step.StepId
				stepsByIdMap[stepId].DepsToClear[parentStepId] = parent
				output[i] = stepsByIdMap[stepId]
			} else {
				return output, fmt.Errorf("invalid dependency specified: %s", parentStepId)
			}
		}
	}

	return output, nil
}

func getNextAvailableStep(stepsByIdSlice []*JobStep) *JobStep {

	possibleNodes := make([]*JobStep, 0)
	for _, step := range stepsByIdSlice {
		if len(step.DepsToClear) == 0 && !step.AllDepsClear {
			possibleNodes = append(possibleNodes, step)
		}
	}

	// We're done.
	if len(possibleNodes) == 0 {
		return nil
	}

	sort.Slice(possibleNodes, func(i, j int) bool {
		if possibleNodes[i].Precedence != possibleNodes[j].Precedence {
			return possibleNodes[i].Precedence > possibleNodes[j].Precedence
		}
		return possibleNodes[i].StepId < possibleNodes[j].StepId
	})

	nodeToReturn := possibleNodes[0]
	for _, node := range stepsByIdSlice {
		newDeps := make(map[string]*JobStep, 0)
		for _, depToClear := range node.DepsToClear {
			if depToClear.StepId == nodeToReturn.StepId {
				continue
			}
			newDeps[depToClear.StepId] = depToClear
		}
		node.DepsToClear = newDeps
	}

	nodeToReturn.AllDepsClear = true

	return nodeToReturn
}
