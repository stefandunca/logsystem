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
}
