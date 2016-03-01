package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type phase struct {
	Name string
	Data *nodeData
}

type setup struct {
	Phases    []*phase
	Variables map[string]interface{}
}

//TODO: Do we need a check the `ContinueOnFailure==true` when `RunParallel==true`. Because if continue is false we will probably exit while another command is still busy in parallel
func (s *setup) Validate() error {
	for _, phase := range s.Phases {
		for _, e := range phase.Data.AdditionalEnvironment {
			if !strings.Contains(e, "=") {
				return fmt.Errorf("Environment string must be in format key=value. Invalid string: '%s'", e)
			}
		}
	}
	return nil
}

func (s *setup) ExpandRepeatingSteps() error {
	for _, phase := range s.Phases {
		err := phase.Data.ExpandRepeatingSteps(s.Variables)
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

func getSinglePhaseFromYaml(yamlBytes []byte, phaseName string) (*phase, error) {
	var err error

	tmpPhases := make(map[string]nodeData)
	if err = yaml.Unmarshal(yamlBytes, &tmpPhases); err != nil {
		return nil, err
	}

	if nodeDat, ok := tmpPhases[phaseName]; !ok {
		return nil, fmt.Errorf("Phase name '%s' was not found", phaseName)
	} else {
		return &phase{Name: phaseName, Data: &nodeDat}, nil
	}
}

func getSetupFromYaml(yamlBytes []byte) (*setup, error) {
	var err error

	tmpPhases := yaml.MapSlice{}
	if err = yaml.Unmarshal(yamlBytes, &tmpPhases); err != nil {
		return nil, err
	}
	tmpPhases = deleteVariablesFromMapSlice(tmpPhases)

	tmpVariables := &struct {
		Variables map[string]interface{}
	}{}
	if err = yaml.Unmarshal(yamlBytes, tmpVariables); err != nil {
		return nil, err
	}

	pointerPhases := []*phase{}
	for _, ph := range tmpPhases {
		phaseName, ok := ph.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Unsupported map KEY type, only string allowed. Key: %#v (%#T)", ph.Key, ph.Key)
		}

		singlePhase, err := getSinglePhaseFromYaml(yamlBytes, phaseName)
		if err != nil {
			return nil, err
		}
		pointerPhases = append(pointerPhases, singlePhase)
	}

	s := &setup{
		pointerPhases,
		tmpVariables.Variables,
	}
	return s, nil
}

func ParseYamlFile(filePath string) (*setup, error) {
	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	s, err := getSetupFromYaml(yamlBytes)
	if err != nil {
		return nil, err
	}

	if err = s.Validate(); err != nil {
		return nil, err
	}

	if err = s.ExpandRepeatingSteps(); err != nil {
		return nil, err
	}

	return s, nil
}
