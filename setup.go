package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"strings"
)

type phasesMap map[string]nodeData

type setup struct {
	Phases    phasesMap
	Variables map[string]string
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

func ParseYamlFile(filePath string) (*setup, error) {
	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	tmpPhases := make(phasesMap)
	if err = yaml.Unmarshal(yamlBytes, &tmpPhases); err != nil {
		return nil, err
	}
	deleteVariablesFromPhasesMap(tmpPhases)

	tmpVariables := &struct {
		Variables map[string]string
	}{}
	if err = yaml.Unmarshal(yamlBytes, tmpVariables); err != nil {
		return nil, err
	}

	s := &setup{
		tmpPhases,
		tmpVariables.Variables,
	}

	if err = s.Validate(); err != nil {
		return nil, err
	}

	return s, nil
}
