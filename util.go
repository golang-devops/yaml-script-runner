package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func splitEnvironKeyValue(pair string) (string, string, error) {
	index := strings.Index(pair, "=")
	if index == -1 {
		return "", "", fmt.Errorf("Cannot find the '=' symbol in environ key-value pair '%s'", pair)
	}
	return pair[0:index], pair[index+1:], nil
}

func appendEnvironment(environ []string, toAppend ...string) ([]string, error) {
	newSlice := []string{}

	for _, e := range environ {
		key1, val1, err := splitEnvironKeyValue(e)
		if err != nil {
			return nil, err
		}

		mustOverwrite := false
		overwriteValue := ""
		for _, a := range toAppend {
			key2, val2, err := splitEnvironKeyValue(a)
			if err != nil {
				return nil, err
			}

			if strings.EqualFold(key1, key2) {
				mustOverwrite = true
				overwriteValue = val2
				break
			}
		}

		if mustOverwrite {
			newSlice = append(newSlice, key1+"="+overwriteValue)
		} else {
			newSlice = append(newSlice, key1+"="+val1)
		}
	}

	for _, a := range toAppend {
		key1, val1, err := splitEnvironKeyValue(a)
		if err != nil {
			return nil, err
		}

		alreadyAdded := false
		for _, n := range newSlice {
			key2, _, err := splitEnvironKeyValue(n)
			if err != nil {
				return nil, err
			}

			if strings.EqualFold(key1, key2) {
				alreadyAdded = true
				break
			}
		}

		if !alreadyAdded {
			newSlice = append(newSlice, key1+"="+val1)
		}
	}

	return newSlice, nil
}

func deleteVariablesFromPhasesMap(m map[string]nodeData) {
	for key, _ := range m {
		if strings.EqualFold(key, "variables") {
			delete(m, key)
			return
		}
	}
}

func replaceVariables(s string, variables map[string]string) string {
	returnStr := s
	for varName, varVal := range variables {
		returnStr = strings.Replace(returnStr, "$"+varName, varVal, -1)
	}
	return returnStr
}

func execTemplateToString(templateString string, data interface{}) (string, error) {
	t, err := template.New("").Parse(templateString)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	err = t.Execute(&doc, data)
	if err != nil {
		return "", err
	}

	return doc.String(), nil
}
