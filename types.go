package main

import (
	"fmt"
	"strings"
)

// Represents the raw unvalidated data coming in from the user input
type InputJobStep struct {
	StepName     string   `yaml:"step"`
	Dependencies []string `yaml:"dependencies"`
	Precedence   int64    `yaml:"precedence"`
}

type InputJob struct {
	Steps []*InputJobStep
}

// Returns whether the user inputted step was valid. If false, error will
// contain the details about the failure and should be handled above.
func (inputStep *InputJobStep) ValidateInputStep() error {

	stepName := inputStep.StepName
	if stepName == "" {
		return fmt.Errorf("no step name provided, invalid")
	}

	stepId := strings.TrimSpace(stepName)
	if stepId == "" {
		return fmt.Errorf("step ID would be empty, invalid name")
	}

	if strings.Contains(stepId, "\n") {
		return fmt.Errorf("newline detected in the StepId")
	}

	return nil
}

// Takes a user input step and returns a fully formed JobStep
func (inputStep *InputJobStep) GetJobStep() *JobStep {

	js := &JobStep{
		StepName:        inputStep.StepName,
		StepId:          strings.TrimSpace(inputStep.StepName),
		Precedence:      inputStep.Precedence,
		DependencyIds:   inputStep.Dependencies,
		DepsToClear:     make(map[string]*JobStep),
		AllDepsClear:    false,
		StepCycleNumber: 0,
	}

	return js
}

// Represents a validated job step with additional data fields for managing
// dependency graph scheduling
type JobStep struct {
	StepName        string              // Represents the original untrimmed Step Name
	StepId          string              // StepName but trimmed of leading and trailing whitespace
	Precedence      int64               // Sorted desc (e.g. Precedence 100 before Precedence 50)
	DependencyIds   []string            // Copies from the Input Dependencies array (represents parentss)
	DepsToClear     map[string]*JobStep // Parent Depdendencies
	AllDepsClear    bool                // Whether all dependencies are clear for this item
	StepCycleNumber int                 // To group steps that can be run at the same time, and need to be further sorted by precedence desc / StepID asc
}

// Convenience function to determine if all parent dependencies have been cleared
func (step *JobStep) AllParentDepsClear() bool {
	for _, parentDep := range step.DepsToClear {
		if !parentDep.AllDepsClear {
			return false
		}
	}

	return true
}
