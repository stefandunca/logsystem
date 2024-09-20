package helpers

import (
	"fmt"
	"logsystem"
	"os"
)

func GetFilePathFromProgArgs() (config logsystem.Config, err error) {
	if len(os.Args) < 2 {
		fmt.Println("Missing config file: \"example <config file>\"")
		return logsystem.Config{}, fmt.Errorf("missing config file")
	}

	configFile := os.Args[1]
	conf, err := logsystem.LoadConfigFromFile(configFile)
	if err != nil {
		fmt.Println("ERROR loading config file:", configFile, "; err:", err)
		os.Exit(2)
	}
	return conf, err
}
