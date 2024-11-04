package main

import (
	"fmt"
	"logsystem"
)

type CountingDriver struct {
	lines int
}

func (d *CountingDriver) Log(data map[logsystem.Param]string) {
	d.lines++
}

func (d *CountingDriver) BeginTx(id logsystem.TxID, attr map[logsystem.Param]string) {
	d.lines++
}

func (d *CountingDriver) EndTx(id logsystem.TxID) {
	d.lines++
}

func (d *CountingDriver) Stop() {
}

func main() {
	m := logsystem.NewManager()

	countDriver := &CountingDriver{}
	serialCouter := logsystem.NewSerialDriver(countDriver)
	m.AddDriver(serialCouter)

	consoleDriver := &logsystem.ConsoleDriver{}
	m.AddDriver(consoleDriver)

	l := logsystem.NewLogger(m)
	l.Info("Hello, world!")
	tl := l.BeginTx(map[logsystem.Param]string{"UserID": "123"})
	tl.Warn("Doing something in TX")
	tl2 := l.BeginTx(map[logsystem.Param]string{"UserID": "456"})
	tl.EndTx()
	l.Debug("Outside TX")
	tl2.Warn("Doing something in TX 2")
	l.Error("Error")
	tl2.EndTx()
	l.Stop()

	defer fmt.Println("Lines logged:", countDriver.lines)
}
