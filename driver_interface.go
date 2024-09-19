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
	driverID() DriverID
	createDriver(config json.RawMessage) DriverInterface
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

type DriverInterface interface {
	log(data map[Param]string)
	beginTx(id TxID)
	endTx(id TxID)
}
