package main

import (
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
const complexWithDependenciesSingleStepLeadingTrailingWhitespaceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "  deploy api gateway "
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

const complexWithDependenciesMultipleStepLeadingTrailingWhitespaceInput string = `
- step: " deploy lambda function"
  dependencies: []
  precedence: 50
- step: " deploy api gateway "
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket   "
  dependencies: []
  precedence: 20
- step: "  enable dns records "
  dependencies: []
  precedence: 200
- step: "     enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const complexWithDependenciesForPrecedenceTestingInput string = `
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

const singleDependenciesWithLeadingWhitespaceInputInput string = `
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

const singleDependenciesWithTrailingWhitespaceInputInput string = `
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

const multipleDependenciesWithTrailingWhitespaceInputInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["  deploy lambda function ", "enable dns records "]
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
  dependencies: ["create bucket  ", "   enable dns records", "deploy database ", " deploy api gateway"]
  precedence: 100
`

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
		{complexWithDependenciesInput, complexWithDependenciesOutput}, // Tests sequencing of parents / children and sorting across 3 dimensions
		{multipleNoDependenciesInput, multipleNoDependenciesOutput},
		{singleDependenciesWithLeadingWhitespaceInputInput, complexWithDependenciesOutput},
		{singleDependenciesWithTrailingWhitespaceInputInput, complexWithDependenciesOutput},
		{multipleDependenciesWithTrailingWhitespaceInputInput, complexWithDependenciesOutput},
		{complexWithDependenciesSingleStepLeadingTrailingWhitespaceInput, complexWithDependenciesOutput},
		{complexWithDependenciesMultipleStepLeadingTrailingWhitespaceInput, complexWithDependenciesOutput},
		{complexWithDependenciesForPrecedenceTestingInput, complexWithDependenciesOutput},
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

const multipleStepsWithSameIdInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "deploy lambda function"
  dependencies: []
  precedence: 40
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

const singleMissingPrecedenceFieldInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: []
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

const multipleMissingPrecedenceFieldInput string = `
- step: "deploy lambda function"
  dependencies: []
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
- step: "enable cdn distribution"
  dependencies: []
  precedence: 100
`

const singleZeroPrecedenceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 0
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

const singleNegativePrecedenceInput string = `
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
  precedence: -20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const singleFloatPrecedenceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 100
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 14.2
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 0
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const singleStringPrecedenceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 100
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: "six"
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 0
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const multipleZeroPrecedenceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 0
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
  precedence: 0
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const multipleNegativePrecedenceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: -50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: -200
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const multipleFloatPrecedenceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 5.0
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 0.25
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 9.5
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const multipleDifferingPrecedenceInput string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 5.0
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 0
- step: "deploy database"
  dependencies: []
  precedence: -50
- step: "create bucket"
  dependencies: []
  precedence: five
- step: "enable dns records"
  dependencies: []
  precedence: nine
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
  precedence: 100
`

const singleInvalidDependency string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns records"]
  precedence: 100
- step: "deploy databaseE"
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

const singleEmptyDependency string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["  ", "enable dns records"]
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

const multipleInvalidDependency string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["deploy lambda function", "enable dns recordss"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "reate bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: []
  precedence: 200
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gatewayy"]
  precedence: 100
`

const multipleEmptyDependency string = `
- step: "deploy lambda function"
  dependencies: []
  precedence: 50
- step: "deploy api gateway"
  dependencies: ["", "enable dns records"]
  precedence: 100
- step: "deploy database"
  dependencies: []
  precedence: 50
- step: "create bucket"
  dependencies: []
  precedence: 20
- step: "enable dns records"
  dependencies: ["    "]
  precedence: 200
- step: "enable cdn distribution"
  dependencies: ["create bucket", "enable dns records", "deploy database", "deploy api gateway"]
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
		{multipleStepsWithSameIdInput},
		{singleMissingPrecedenceFieldInput},
		{multipleMissingPrecedenceFieldInput},
		{singleZeroPrecedenceInput},
		{singleNegativePrecedenceInput},
		{singleFloatPrecedenceInput},
		{singleStringPrecedenceInput},
		{multipleZeroPrecedenceInput},
		{multipleNegativePrecedenceInput},
		{multipleFloatPrecedenceInput},
		{multipleDifferingPrecedenceInput},
		{singleInvalidDependency},
		{singleEmptyDependency},
		{multipleInvalidDependency},
		{multipleEmptyDependency},
	}

	for i, testCase := range tests {
		_, outputErr := ProcessUserJob(testCase.invalidYamlInput)
		if outputErr == nil {
			t.Errorf("test %d: should have received error, did not", i)
		}
	}
}

func TestAllStepsMustBeUsedExactlyOnce(t *testing.T) {

	output, outputErr := ProcessUserJob(complexWithDependenciesInput)
	if outputErr != nil {
		t.Errorf("received unexpected error: %s", outputErr.Error())
	}

	counter := make(map[string]int, len(output))
	for _, outputStr := range output {
		if _, stringExists := counter[outputStr]; stringExists {
			t.Errorf("duplicate value when processing")
			counter[outputStr]++
		} else {
			counter[outputStr] = 1
		}
	}

	for k, v := range counter {
		if v != 1 {
			t.Errorf("duplicate values for %s, %d", k, v)
		}
	}
}

/**
The rules:

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
