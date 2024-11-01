package main

import (
	"fmt"
	"logsystem"
	"logsystem/example/helpers"
	"os"
)

func main() {
	conf, err := helpers.GetFilePathFromProgArgs()
	if err != nil {
		os.Exit(1)
	}

	m, err := logsystem.CreateLogManagerWithConfig(
		[]logsystem.DriverFactoryInterface{
			&logsystem.ConsoleDriverFactory{},
			&logsystem.FileDriverFactory{},
			&logsystem.DBDriverFactory{},
		},
		conf,
	)
	if err != nil {
		if err != logsystem.ErrorSomeDriversFailed {
			os.Exit(1)
		}
		fmt.Println("Some drivers failed to initialize; check logs for details")
	}

	l := logsystem.NewLogger(m)
	defer l.Stop()
	l.Info("Hello, world!")
	tl := l.BeginTx(map[logsystem.Param]string{"UserID": "123"})
	tl.Warn("Doing something in TX")
	tl2 := l.BeginTx(map[logsystem.Param]string{"UserID": "456"})
	defer tl2.EndTx()
	tl.EndTx()
	l.Debug("Outside TX")
	tl2.Warn("Doing something in TX 2")
	l.Error("Error")
}
