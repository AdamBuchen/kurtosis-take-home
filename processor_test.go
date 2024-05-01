package main

import (
	"fmt"
	"testing"
)

const oneStepInput string = `
- step: "prepare database"
  dependencies: []
  precedence: 50
`

var oneStepOutput = []string{
	"prepare database",
}

const basicWithDependenciesInput string = `
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

var basicWithDependenciesOutput = []string{
	"prepare database",
	"create user 1",
	"create user 2",
	"create user 4",
	"create user 3",
}

const complexWithDependenciesInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

var complexWithDependenciesOutput = []string{
	"enable dns records",
	"deploy database",
	"deploy lambda function",
	"create bucket",
	"deploy api gateway",
	"enable cdn distribution",
}

const multipleNoDependenciesInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

var multipleNoDependenciesOutput = []string{
	"enable dns records",
	"deploy api gateway",
	"enable cdn distribution",
	"deploy database",
	"deploy lambda function",
	"create bucket",
}

// Convenience function to check whether the values of two string arrays are equal
func areStringArraysEqual(a, b []string) bool {

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestCorrectlyProcessesValidInputs(t *testing.T) {
	var tests = []struct {
		validYamlInput string
		correctOutput  []string
	}{
		{oneStepInput, oneStepOutput},
		{basicWithDependenciesInput, basicWithDependenciesOutput},
		{complexWithDependenciesInput, complexWithDependenciesOutput},
		{multipleNoDependenciesInput, multipleNoDependenciesOutput},
	}

	for i, testCase := range tests {
		output, outputErr := ProcessUserJob(testCase.validYamlInput)
		if outputErr != nil {
			t.Errorf("test %d: expected success, got error: %s", i, outputErr.Error())
		}

		if !areStringArraysEqual(output, testCase.correctOutput) {
			t.Errorf("test %d: did not get equal string arrays, got: %v", i, output)
		}
	}
}

const circularDependenciesInput string = `
- step: "deploy lambda function"
  dependencies: ["deploy api gateway"]
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const selfDependencyInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy api gateway"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`
const nonYamlStringInput string = `
The quuick brown fox jumps over the lazy dog.
`
const nonYamlJsonInput string = `
{
	"key": "extremely unexpected value"
}
`
const emptyStringInput string = ``

const missingSingleStepFieldInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

const missingMultipleStepFieldInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- dependencies: []
  precedence: 100
`

const singleStepWithEmptyInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: ""
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

const multipleStepsWithEmptyInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step:
  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: ""
  dependencies: []
  precedence: 100
`

const singleStepWithWhitespaceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "    "
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`
const singleStepWithNewlineInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: |
         
     
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

const multipleStepsWithWhitespaceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "   "
  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "     "
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

const multipleStepsWithNewlineInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: |
      


  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: |
    
     



  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

const singleStepWithNewlineInNameInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: []
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: |
    This is the first thing
    This is another thing
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

func TestCorrectlyRejectsInvalidInputs(t *testing.T) {
	var tests = []struct {
		invalidYamlInput string
	}{
		{circularDependenciesInput},
		{selfDependencyInput},
		{nonYamlStringInput},
		{nonYamlJsonInput},
		{emptyStringInput},
		{missingSingleStepFieldInput},
		{missingMultipleStepFieldInput},
		{singleStepWithEmptyInput},
		{multipleStepsWithEmptyInput},
		{singleStepWithWhitespaceInput},
		{singleStepWithNewlineInput},
		{multipleStepsWithWhitespaceInput},
		{multipleStepsWithNewlineInput},
		{singleStepWithNewlineInNameInput},
	}

	for i, testCase := range tests {
		_, outputErr := ProcessUserJob(testCase.invalidYamlInput)
		if outputErr == nil {
			t.Errorf("test %d: should have received error, did not", i)
		} else {
			fmt.Println(outputErr.Error())
		}
	}

}

func TestMultipleStepsWithTheSameIdAreNotAllowed(t *testing.T) {

}

func TestEachStepMustHaveAPrecdenceField(t *testing.T) {

}

func TestPrecedenceMustBeAPositiveNonzeroInteger(t *testing.T) {

}

func TestDependencyIdsAreValuesOfDependenciesFieldWithLeadingAndTrailingWhitespaceRemoved(t *testing.T) {

}

func TestAnEmptyOrWhitespaceDependencyIsNotAllowed(t *testing.T) {

}

func TestDependenciesOnNonexistentStepIdsAreNotAllowed(t *testing.T) {

}

func TestOutputOrderingIsNewlineSeparatedListOfStepIdsWithNoLeadingOrTrailingWhitespace(t *testing.T) {

}

func TestOutputOrderingIsAlwaysTerminatedByANewline(t *testing.T) {

}

func TestAllStepsMustBeUsedExactlyOnce(t *testing.T) {

}

func TestIdOfAStepIsValueOfStepFieldWithLeadingAndTrailingWhitespaceRemoved(t *testing.T) {

}

func TestAStepsDependenciesMustComeBeforeItInOutputOrdering(t *testing.T) {

}

func TestHigherPrecedenceStepsMustComeBeforeLowerPrecedenceStepsWhenBothAreAvailableForRunning(t *testing.T) {

}

func TestWhenTwoReadyToRunStepsHaveTheSamePrecedenceUseLexicographicalOrdering(t *testing.T) {

}

/**
A job is a list of steps defined in YAML
A job must have at least one step
Each step must have a step field
The ID of a step is the value of the step field, with the leading and trailing whitespace removed
Empty or all-whitespace step IDs are not allowed
Step IDs with newline characters are not allowed
Multiple steps with the same ID are not allowed
Each step must have a precedence field
The precedence of a step is the value of the precedence field
Precedence must be a positive nonzero integer
Each step may (but is not required to) have a dependencies field containing an array of step IDs
A step without the dependencies key is assumed to have no dependencies
The dependency IDs are the values of the dependencies, field with leading and trailing whitespace removed
An empty or whitespace dependency ID is not allowed
Dependencies on nonexistent step IDs are not allowed

An output ordering is a newline-separated list of step IDs (no leading or trailing whitespace)
An output ordering is always terminated by a newline
All steps in the job must be used exactly once
A step's dependencies must come before it in the output ordering
Higher-precedence steps must come before lower-precedence steps when both are available for running
When two ready-to-run steps have the same precedence, lexicographical ordering is used: step A comes before step B, etc.

**/
