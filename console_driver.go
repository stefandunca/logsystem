package logsystem

import (
	"encoding/json"
	"fmt"
)

const ConsoleDriverID = "console"

type consoleConfig struct {
	UserReadableTime bool `json:"userReadableTime"`
}

// ConsoleDriverFactory implements DriverFactoryInterface
type ConsoleDriverFactory struct {
}

func (f *ConsoleDriverFactory) driverID() DriverID {
	return DriverID(ConsoleDriverID)
}

func (f *ConsoleDriverFactory) createDriver(config json.RawMessage) DriverInterface {
	var consoleConfig consoleConfig
	err := json.Unmarshal(config, &consoleConfig)
	if err != nil {
		fmt.Printf("Failed to unmarshal console driver config: %v\n", err)
		return nil
	}

	return &ConsoleDriver{
		config: consoleConfig,
	}
}

// ConsoleDriver implements DriverInterface
type ConsoleDriver struct {
	config consoleConfig
}

func (d *ConsoleDriver) log(data map[Param]string) {
	line := formatLine(data, d.config.UserReadableTime)
	fmt.Println(line)
}

func (d *ConsoleDriver) beginTx(id TxID, attr map[Param]string) {
	txData := make(map[Param]string)
	txData[TxIDParam] = id.String()
	message := fmt.Sprintf("TX Begin; Params: %v", attr)
	txData[MessageParam] = message
	txData[LevelParam] = string(Info)
	d.log(txData)
}

func (d *ConsoleDriver) endTx(id TxID) {
	txData := make(map[Param]string)
	txData[TxIDParam] = id.String()
	txData[MessageParam] = "TX End"
	txData[LevelParam] = string(Info)
	d.log(txData)
}

func (d *ConsoleDriver) stop() {
}
