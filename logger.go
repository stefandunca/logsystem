package logsystem

import (
	"strconv"
	"time"
)

type Logger struct {
	mgr *DriverManager
}

func NewLogger(conf Config) *Logger {
	factories := []DriverFactoryInterface{
		&ConsoleDriverFactory{},
	}
	return NewLoggerWithDrivers(conf, factories)
}

func NewLoggerWithDrivers(conf Config, factories []DriverFactoryInterface) *Logger {
	return &Logger{
		mgr: NewManager(factories, conf),
	}
}

func (l *Logger) Info(message string) {
	l.logBasic(message, Info)
}

func (l *Logger) logBasic(message string, level LogLevel) {
	data := map[Param]string{
		MessageParam:   message,
		TimeParam:      strconv.FormatInt(time.Now().Unix(), 10),
		LevelParam:     string(level),
		ComponentParam: "component",
	}
	l.mgr.log(data)
}
