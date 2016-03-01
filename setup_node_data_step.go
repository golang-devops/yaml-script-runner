package main

import (
	"github.com/golang-devops/parsecommand"
)

type nodeDataStep string

func (n nodeDataStep) SplitAndReplaceVariables(stringVariables map[string]string, allVariables map[string]interface{}) ([]string, error) {
	step := replaceVariables(string(n), stringVariables)

	splittedStep, err := parsecommand.Parse(step)
	if err != nil {
		return nil, err
	}

	return splittedStep, nil
}
