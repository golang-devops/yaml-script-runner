package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var BuildSha1Version string

func checkError(err error, prefix string) {
	if err != nil {
		panic(fmt.Sprintf("%s %s", prefix, err.Error()))
	}
}

func cleanFeedbackLine(s string) (out string) {
	out = s
	out = strings.TrimSpace(out)
	out = strings.Replace(out, "\r", "", -1)
	out = strings.Replace(out, "\n", "\\n", -1)
	return
}

func runCommand(wg *sync.WaitGroup, resultCollector *ResultCollector, commandIndex int, phase *nodeData, cmd *exec.Cmd) {
	if wg != nil {
		defer wg.Done()
	}

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("%s", r)
			if !phase.ContinueOnFailure {
				logger.Fatallnf(errMsg)
			} else {
				logger.Errorlnf(errMsg)
			}

			resultCollector.AppendResult(&result{Cmd: cmd, Successful: false})
		} else {
			resultCollector.AppendResult(&result{Cmd: cmd, Successful: true})
		}
	}()

	stdout, err := cmd.StdoutPipe()
	checkError(err, "cmd.StdoutPipe")

	stderr, err := cmd.StderrPipe()
	checkError(err, "cmd.StderrPipe")

	outLines := []string{}
	stdoutScanner := bufio.NewScanner(stdout)
	go func() {
		for stdoutScanner.Scan() {
			txt := cleanFeedbackLine(stdoutScanner.Text())
			outLines = append(outLines, txt)
			logger.Tracelnf("COMMAND %d STDOUT: %s", commandIndex, txt)
		}
	}()

	errLines := []string{}
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			txt := cleanFeedbackLine(stderrScanner.Text())
			errLines = append(errLines, txt)
			logger.Warninglnf("COMMAND %d STDERR: %s", commandIndex, txt)
		}
	}()

	err = cmd.Start()
	checkError(err, "cmd.Start")

	err = cmd.Wait()
	if err != nil {
		newErr := fmt.Errorf("%s - lines: %s", err.Error(), strings.Join(errLines, "\\n"))
		outCombined := cleanFeedbackLine(strings.Join(outLines, "\n"))

		errMsg := fmt.Sprintf("FAILED COMMAND %d. Continue: %t. ERROR: %s. OUT: %s", commandIndex, phase.ContinueOnFailure, newErr.Error(), outCombined)
		panic(errMsg)
	}
}

func runPhase(setup *setup, phaseName string, phase *nodeData) {
	var wg sync.WaitGroup
	resultCollector := &ResultCollector{}

	cmds, err := phase.GetExecCommandsFromSteps(setup.StringVariablesOnly(), setup.Variables)
	if err != nil {
		logger.Fatallnf(err.Error())
	}

	if phase.RunParallel {
		wg.Add(len(cmds))
	}

	logger.Infolnf("Running step %s", phaseName)
	for ind, c := range cmds {
		logger.Tracelnf(strings.TrimSpace(fmt.Sprintf(`INDEX %d = %s`, ind, execCommandToDisplayString(c))))

		var wgToUse *sync.WaitGroup = nil
		if phase.RunParallel {
			wgToUse = &wg
			go runCommand(wgToUse, resultCollector, ind, phase, c)
		} else {
			runCommand(wgToUse, resultCollector, ind, phase, c)
		}
	}

	if phase.RunParallel {
		wg.Wait()
	}

	if resultCollector.FailedCount() == 0 {
		logger.Infolnf("There were no failures - %d/%d successful", resultCollector.SuccessCount(), resultCollector.TotalCount())
	} else {
		logger.Errorlnf("There were %d/%d failures: ", resultCollector.FailedCount(), resultCollector.TotalCount())
		for _, failure := range resultCollector.FailedDisplayList() {
			logger.Errorlnf("- %s", failure)
		}
	}
}

func main() {
	logger.Infolnf("Running version '%s'", BuildSha1Version)

	if len(os.Args) < 2 {
		logger.Fatallnf("The first command-line argument must be the YAML file path.")
	}

	yamlFilePath := os.Args[1]

	setup, err := ParseYamlFile(yamlFilePath)
	if err != nil {
		logger.Fatallnf(err.Error())
	}

	for _, phase := range setup.Phases {
		runPhase(setup, phase.Name, phase.Data)
	}
}
