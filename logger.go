package logsystem

import (
	"strconv"
	"time"
)

type Logger struct {
	mgr *DriverManager
}

type TxLogger struct {
	logger    *Logger
	txID      TxID
	component string
}

func NewLogger(conf Config) *Logger {
	factories := []DriverFactoryInterface{
		&ConsoleDriverFactory{},
		&FileDriverFactory{},
		&DBDriverFactory{},
	}
	return NewLoggerWithDrivers(conf, factories)
}

func NewLoggerWithDrivers(conf Config, factories []DriverFactoryInterface) *Logger {
	return &Logger{
		mgr: NewManager(factories, conf),
	}
}

func (l *Logger) Stop() {
	l.mgr.stop()
}

func (l *Logger) Info(message string) {
	l.logBasic(message, Info)
}

func (l *Logger) Debug(message string) {
	l.logBasic(message, Debug)
}

func (l *Logger) Warn(message string) {
	l.logBasic(message, Warn)
}

func (l *Logger) Error(message string) {
	l.logBasic(message, Error)
}

func (l *Logger) BeginTx(attr map[Param]string) TxLogger {
	return l.BeginTxWithComponent("", attr)
}

func (l *Logger) BeginTxWithComponent(component string, attr map[Param]string) TxLogger {
	txID := l.mgr.beginTx(attr)
	tx := TxLogger{
		logger:    l,
		txID:      txID,
		component: component,
	}
	return tx
}

func (tl TxLogger) Info(message string) {
	tl.logAttrib(message, Info)
}

func (tl TxLogger) Debug(message string) {
	tl.logAttrib(message, Debug)
}

func (tl TxLogger) Warn(message string) {
	tl.logAttrib(message, Warn)
}

func (tl TxLogger) Error(message string) {
	tl.logAttrib(message, Error)
}

func (tl TxLogger) logAttrib(message string, level LogLevel) {
	extra := map[Param]string{
		TxIDParam: tl.txID.String(),
	}
	if tl.component != "" {
		extra[ComponentParam] = tl.component
	}
	tl.logger.logAttrib(message, level, extra)
}

func (tl TxLogger) EndTx() {
	tl.logger.mgr.endTx(tl.txID)
}

func (l *Logger) logBasic(message string, level LogLevel) {
	l.logAttrib(message, level, nil)
}

func (l *Logger) logAttrib(message string, level LogLevel, attributes map[Param]string) {
	data := map[Param]string{
		MessageParam: message,
		TimeParam:    strconv.FormatInt(time.Now().Unix(), 10),
		LevelParam:   string(level),
	}
	for k, v := range attributes {
		data[k] = v
	}
	l.mgr.log(data)
}
