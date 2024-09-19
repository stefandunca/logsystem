package logsystem

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
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

func formatLine(data map[Param]string, userFriendly bool) string {
	var timestamp int64 = 0
	if val, ok := data[TimeParam]; ok {
		parsedTime, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			timestamp = parsedTime
		}
	}

	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	formattedTime := ""
	if userFriendly {
		t := time.Unix(timestamp, 0)
		formattedTime = t.Format("[2006-01-02 15:04:05] ")
	} else {
		formattedTime = fmt.Sprintf("[%-10d] ", timestamp)
	}

	level := ""
	if val, ok := data[LevelParam]; ok {
		level = strings.ToUpper(val)
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

	return fmt.Sprintf("%s%-5s %s%s", formattedTime, level, message, optional)
}
