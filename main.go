package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func runCommand(wg *sync.WaitGroup, commandIndex int, phase *nodeData, cmd *exec.Cmd) {
	if wg != nil {
		defer wg.Done()
	}

	commandName := strings.TrimSpace(fmt.Sprintf(`(INDEX %d) "%s" %+v`, commandIndex, cmd.Path, cmd.Args))

	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("ERROR (continue=%t): %s. OUT: %s. COMMAND: %s\n", phase.ContinueOnFailure, err.Error(), string(out), commandName)
		if !phase.ContinueOnFailure {
			logger.Fatallnf(errMsg)
		} else {
			logger.Errorlnf(errMsg)
			return
		}
	}

	logger.PrintCommandOutput(commandName, string(out))
}

func runPhase(setup *setup, phaseName string, phase *nodeData) {
	var wg sync.WaitGroup

	cmds, err := phase.GetExecCommandsFromSteps(setup.StringVariablesOnly(), setup.Variables)
	if err != nil {
		logger.Fatallnf(err.Error())
	}

	if phase.RunParallel {
		wg.Add(len(cmds))
	}

	logger.Infolnf("Running step %s", phaseName)
	for ind, c := range cmds {
		var wgToUse *sync.WaitGroup = nil
		if phase.RunParallel {
			wgToUse = &wg
			go runCommand(wgToUse, ind, phase, c)
		} else {
			runCommand(wgToUse, ind, phase, c)
		}
	}

	if phase.RunParallel {
		wg.Wait()
	}
}

func main() {
	if len(os.Args) < 2 {
		logger.Fatallnf("The first command-line argument must be the YAML file path.")
	}

	yamlFilePath := os.Args[1]

	setup, err := ParseYamlFile(yamlFilePath)
	if err != nil {
		logger.Fatallnf(err.Error())
	}

	for name, phase := range setup.Phases {
		runPhase(setup, name, phase)
	}
}
