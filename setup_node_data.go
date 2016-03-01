package main

import (
	"os"
	"os/exec"
)

type nodeData struct {
	ContinueOnFailure     bool `json:"continue_on_failure"`
	InheritEnvironment    bool `json:"inherit_environment"`
	AdditionalEnvironment []string
	RunParallel           bool `json:"run_parallel"`
	Executor              []string
	Steps                 []nodeDataStep
}

func (n *nodeData) GetExecCommandsFromSteps(variables map[string]string) ([]*exec.Cmd, error) {
	cmds := []*exec.Cmd{}

	for _, step := range n.Steps {
		splittedStep, err := step.SplitAndReplaceVariables(variables)
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
