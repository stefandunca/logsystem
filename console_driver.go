package logsystem

import (
	"encoding/json"
	"fmt"
)

// ConsoleDriverFactory implements DriverFactoryInterface
type ConsoleDriverFactory struct {
}

func (f *ConsoleDriverFactory) driverID() DriverID {
	return DriverID("_console")
}

func (f *ConsoleDriverFactory) createDriver(config json.RawMessage) DriverInterface {
	return &ConsoleDriver{}
}

// ConsoleDriver implements DriverInterface
type ConsoleDriver struct {
}

func (d *ConsoleDriver) log(data map[Param]string) {
	line := formatLine(data)
	fmt.Println(line)
}

func (d *ConsoleDriver) beginTx(id TxID) {

}

func (d *ConsoleDriver) endTx(txID TxID) {

}

func formatLine(data map[Param]string) string {
	time := ""
	if val, ok := data[TimeParam]; ok {
		time = val
	}

	level := ""
	if val, ok := data[LevelParam]; ok {
		level = val
	}

	message := ""
	if val, ok := data[MessageParam]; ok {
		message = val
	}

	component := ""
	if val, ok := data[ComponentParam]; ok {
		component = val
	}
	txID := ""
	if val, ok := data[TxIDParam]; ok {
		txID = val
	}
	optional := ""
	if txID != "" || component != "" {
		if component != "" {
			optional = fmt.Sprintf("; Comp=[%s]", component)
		}
		if txID != "" {
			optional = fmt.Sprintf("%s; TxID=[%s]", optional, txID)
		}
	}

	return fmt.Sprintf("[%-10s] %-10s %s%s", time, level, message, optional)
}
