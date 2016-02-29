package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("The first command-line argument must be the YAML file path.")
	}

	yamlFilePath := os.Args[1]

	setup, err := ParseYamlFile(yamlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	cmds, err := setup.GetExecCommandsFromSteps()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cmds {
		out, err := c.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("ERROR (continue=%t): %s. OUT: %s\n", setup.ContinueIfStepFailed, err.Error(), string(out))
			if !setup.ContinueIfStepFailed {
				log.Fatalln(errMsg)
			} else {
				log.Println(errMsg)
				continue
			}
		}

		fmt.Println(string(out))
	}
}
