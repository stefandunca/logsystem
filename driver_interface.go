package logsystem

import (
	"encoding/json"
	"strconv"
)

type DriverID string
type TxID int64

func (id TxID) String() string {
	return strconv.Itoa(int(id))
}

type DriverFactoryInterface interface {
	DriverID() DriverID
	CreateDriver(config json.RawMessage) (DriverInterface, error)
}

type Param string

const (
	MessageParam   Param = "message"
	TimeParam      Param = "time"  // Unix timestamp
	LevelParam     Param = "level" // LogLevel
	ComponentParam Param = "component"
	TxIDParam      Param = "txID"
)

type LogLevel string

const (
	Info  LogLevel = "info"
	Debug LogLevel = "debug"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

// DriverInterface interface won't be called if the driver is not created successfully, therefore no need to handle creation errors
type DriverInterface interface {
	Log(data map[Param]string)
	BeginTx(id TxID, attr map[Param]string)
	EndTx(id TxID)

	// shutdown the driver
	Stop()
}
