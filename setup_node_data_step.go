package main

import (
	"github.com/golang-devops/parsecommand"
)

type nodeDataStep string

//TODO: Is the best practice to replace variables after splitting or before?
func (n nodeDataStep) SplitAndReplaceVariables(variables map[string]string) ([]string, error) {
	preReplacedVars := replaceVariables(string(n), variables)
	splittedStep, err := parsecommand.Parse(preReplacedVars)
	if err != nil {
		return nil, err
	}

	/*for i, _ := range splittedStep {
		splittedStep[i] = replaceVariables(splittedStep[i], variables)
	}*/

	return splittedStep, nil
}
