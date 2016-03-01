package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type nodeData struct {
	ContinueOnFailure     bool     `json:"continue_on_failure"`
	InheritEnvironment    bool     `json:"inherit_environment"`
	AdditionalEnvironment []string `json:"additional_environment"`
	RunParallel           bool     `json:"run_parallel"`
	Executor              []string
	Steps                 []nodeDataStep
}

func (n *nodeData) ExpandRepeatingSteps(variables map[string]interface{}) error {
	expandedSteps := []nodeDataStep{}

	for _, step := range n.Steps {
		repeatPrefix := "repeat::"
		if strings.HasPrefix(string(step), repeatPrefix) {
			varNameAndStepLine := strings.TrimPrefix(string(step), repeatPrefix)

			indexFirstSpace := strings.Index(varNameAndStepLine, " ")
			if indexFirstSpace == -1 {
				return fmt.Errorf("Unable to find a space after the repeat prefix '%s' in step line '%s'", repeatPrefix, string(step))
			}

			varName := strings.TrimSpace(varNameAndStepLine[0:indexFirstSpace])
			if varName == "" {
				return fmt.Errorf("The variable name is blank in the repeat prefixed step line '%s'", string(step))
			}

			var varSlice []interface{}
			if val, ok := variables[varName]; !ok {
				return fmt.Errorf("The variable '%s' is not found in variables list %#v", varName, variables)
			} else if varSlice, ok = val.([]interface{}); !ok {
				return fmt.Errorf("The variable '%s' with value '%#v' cannot be type-casted to string slice", varName, val)
			}

			repeatingStep := strings.TrimSpace(varNameAndStepLine[indexFirstSpace:])
			if repeatingStep == "" {
				return fmt.Errorf("The remaining step line after the repeat prefix '%s' is empty or white space in step line '%s'", repeatPrefix, string(step))
			}

			for _, v := range varSlice {
				expandedStep, err := execTemplateToString(repeatingStep, v)
				if err != nil {
					return err
				}
				expandedSteps = append(expandedSteps, nodeDataStep(strings.TrimSpace(expandedStep)))
			}
		} else {
			expandedSteps = append(expandedSteps, step)
		}
	}
	n.Steps = expandedSteps

	return nil
}

func (n *nodeData) GetExecCommandsFromSteps(stringVariables map[string]string, allVariables map[string]interface{}) ([]*exec.Cmd, error) {
	cmds := []*exec.Cmd{}

	for _, step := range n.Steps {
		splittedStep, err := step.SplitAndReplaceVariables(stringVariables, allVariables)
		if err != nil {
			return nil, err
		}

		allArgs := n.Executor
		allArgs = append(allArgs, splittedStep...)

		exe := allArgs[0]
		args := []string{}
		if len(allArgs) > 1 {
			args = allArgs[1:]
		}

		c := exec.Command(exe, args...)

		if n.InheritEnvironment {
			c.Env, err = appendEnvironment(c.Env, os.Environ()...)
			if err != nil {
				return nil, err
			}
		}

		for _, e := range n.AdditionalEnvironment {
			c.Env, err = appendEnvironment(c.Env, e)
			if err != nil {
				return nil, err
			}
		}

		cmds = append(cmds, c)
	}

	return cmds, nil
}
