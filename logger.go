package main

import (
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"sync"
)

var logger = NewLogger(os.Stdout /*ioutil.Discard*/, os.Stdout, os.Stdout, os.Stderr)

func NewLogger(traceHandle, infoHandle, warningHandle, errorHandle io.Writer) *Logger {
	return &Logger{
		colorTraceSprintfFunc:   color.New().SprintfFunc(),
		colorInfoSprintfFunc:    color.New(color.FgHiGreen).SprintfFunc(),
		colorWarningSprintfFunc: color.New(color.FgHiYellow).SprintfFunc(),
		colorErrSprintfFunc:     color.New(color.FgRed).SprintfFunc(),

		trace:   log.New(traceHandle, "[T] ", log.Ldate|log.Ltime|log.Lshortfile),
		info:    log.New(infoHandle, "[I] ", log.Ldate|log.Ltime|log.Lshortfile),
		warning: log.New(warningHandle, "[W] ", log.Ldate|log.Ltime|log.Lshortfile),
		err:     log.New(errorHandle, "[E] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

type Logger struct {
	sync.RWMutex

	colorTraceSprintfFunc   func(format string, args ...interface{}) string
	colorInfoSprintfFunc    func(format string, args ...interface{}) string
	colorWarningSprintfFunc func(format string, args ...interface{}) string
	colorErrSprintfFunc     func(format string, args ...interface{}) string

	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
}

func (l *Logger) Tracelnf(format string, args ...interface{}) {
	l.trace.Println(l.colorTraceSprintfFunc(format+"\n", args...))
}
func (l *Logger) Infolnf(format string, args ...interface{}) {
	l.info.Println(l.colorInfoSprintfFunc(format+"\n", args...))
}
func (l *Logger) Warninglnf(format string, args ...interface{}) {
	l.warning.Println(l.colorWarningSprintfFunc(format+"\n", args...))
}
func (l *Logger) Errorlnf(format string, args ...interface{}) {
	l.err.Println(l.colorErrSprintfFunc(format+"\n", args...))
}
func (l *Logger) Fatallnf(format string, args ...interface{}) {
	l.err.Fatalln(l.colorErrSprintfFunc(format+"\n", args...))
}

func (l *Logger) PrintCommandOutput(commandName string, output string) {
	l.Lock()
	defer l.Unlock()

	l.Infolnf("[BEGIN OUTPUT of %s]", commandName)
	color.New(color.FgHiCyan).Println(string(output))
	l.Infolnf("[END OUTPUT of %s]", commandName)
}
