package main

import (
	"fmt"
	"logsystem"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing config file: \"example <config file>\"")
		os.Exit(1)
	}

	configFile := os.Args[1]
	conf, err := logsystem.LoadConfigFromFile(configFile)
	if err != nil {
		fmt.Println("ERROR loading config file:", configFile, "; err:", err)
		os.Exit(2)
	}

	l := logsystem.NewLogger(conf)
	l.Info("Hello, world!")
	tl := l.BeginTx(map[logsystem.Param]string{"UserID": "123"})
	tl.Warn("Doing something in TX")
	tl2 := l.BeginTx(map[logsystem.Param]string{"UserID": "456"})
	defer tl.EndTx()
	tl.EndTx()
	l.Debug("Outside TX")
	tl2.Warn("Doing something in TX 2")
	l.Error("Error")
}
