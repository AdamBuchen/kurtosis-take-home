package main

import "fmt"

// Return a map which is keyed by the id of the step
func getStepsByIdMap(inputSteps []*InputJobStep) (map[string]*JobStep, error) {

	output := make(map[string]*JobStep)
	for _, inputStep := range inputSteps {
		validationErr := inputStep.ValidateInputStep()
		if validationErr != nil {
			return output, fmt.Errorf("validation error received: " + validationErr.Error())
		}

		jobStep := inputStep.GetJobStep()
		output[jobStep.StepId] = jobStep
	}

	//Now we loop over our output to make sure that there are no invalid dependencies
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

// Map keyed by parent, each containing an array of direct children steps
func getDepsByParent(stepsByIdMap map[string]*JobStep) (map[string][]*JobStep, error) {

	output := make(map[string][]*JobStep)
	for stepId, _ := range stepsByIdMap {
		output[stepId] = make([]*JobStep, 0)
	}

	for _, step := range stepsByIdMap {
		for parentId, _ := range step.DepsToClear {
			output[parentId] = append(output[parentId], step)
		}
	}

	return output, nil
}

// 4. Take the stepsByIdMap and get a depsByParent map
//depsByParent := make(map[string][]*JobStep)

// A recursive function that checks whther in the current cycle we can yet
// process this dependency. Results in Group Number being set and AllDepsClear
// being set to true for any children that can be processed.
func processChildItem(step *JobStep, depsByParent map[string][]*JobStep, currentGroupNumber int) error {

	if step.StepGroupNumber > 0 && currentGroupNumber > step.StepGroupNumber {
		return fmt.Errorf("cyclical relationship detected: %s", step.StepId)
	}

	// This one was already processed, let's not get stuck in a loop
	if step.AllDepsClear || step.StepGroupNumber > 0 {
		return nil
	}

	if step.AllParentDepsClear() {
		step.AllDepsClear = true
		step.StepGroupNumber = currentGroupNumber
		if _, ok := depsByParent[step.StepId]; ok {
			nextGroupNumber := currentGroupNumber + 1
			for _, childStep := range depsByParent[step.StepId] {
				err := processChildItem(childStep, depsByParent, nextGroupNumber)
				if err != nil {
					return err // Just pass it up the chain
				}
			}
		}
	}

	return nil
}
