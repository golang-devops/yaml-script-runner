package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"strings"
)

type phasesMap map[string]*nodeData

type setup struct {
	Phases    phasesMap
	Variables map[string]interface{}
}

//TODO: Do we need a check the `ContinueOnFailure==true` when `RunParallel==true`. Because if continue is false we will probably exit while another command is still busy in parallel
func (s *setup) Validate() error {
	for _, node := range s.Phases {
		for _, e := range node.AdditionalEnvironment {
			if !strings.Contains(e, "=") {
				return fmt.Errorf("Environment string must be in format key=value. Invalid string: '%s'", e)
			}
		}
	}
	return nil
}

func (s *setup) ExpandRepeatingSteps() error {
	for _, node := range s.Phases {
		err := node.ExpandRepeatingSteps(s.Variables)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *setup) StringVariablesOnly() map[string]string {
	m := make(map[string]string)
	for k, v := range s.Variables {
		if str, ok := v.(string); ok {
			m[k] = str
		}
	}
	return m
}

func ParseYamlFile(filePath string) (*setup, error) {
	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	tmpPhases := make(map[string]nodeData)
	if err = yaml.Unmarshal(yamlBytes, &tmpPhases); err != nil {
		return nil, err
	}
	deleteVariablesFromPhasesMap(tmpPhases)

	tmpVariables := &struct {
		Variables map[string]interface{}
	}{}
	if err = yaml.Unmarshal(yamlBytes, tmpVariables); err != nil {
		return nil, err
	}

	pointerPhasesMap := make(phasesMap)
	for key, phase := range tmpPhases {
		pointerPhasesMap[key] = &phase
	}

	s := &setup{
		pointerPhasesMap,
		tmpVariables.Variables,
	}

	if err = s.Validate(); err != nil {
		return nil, err
	}

	if err = s.ExpandRepeatingSteps(); err != nil {
		return nil, err
	}

	return s, nil
}
