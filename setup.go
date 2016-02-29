package main

import (
	"github.com/ghodss/yaml"
	"github.com/golang-devops/parsecommand"
	"io/ioutil"
	"os/exec"
)

type setup struct {
	ContinueIfStepFailed bool `json:"continue_if_step_failed"`
	Executor             []string
	Steps                []string
}

func (s *setup) GetExecCommandsFromSteps() ([]*exec.Cmd, error) {
	cmds := []*exec.Cmd{}

	for _, step := range s.Steps {
		splittedStep, err := parsecommand.Parse(step)
		if err != nil {
			return nil, err
		}

		allArgs := s.Executor
		allArgs = append(allArgs, splittedStep...)

		exe := allArgs[0]
		args := []string{}
		if len(allArgs) > 1 {
			args = allArgs[1:]
		}

		cmds = append(cmds, exec.Command(exe, args...))
	}

	return cmds, nil
}

func ParseYamlFile(filePath string) (*setup, error) {
	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	s := &setup{}
	err = yaml.Unmarshal(yamlBytes, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
