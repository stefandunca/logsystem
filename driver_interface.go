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
	MessageParam   Param = "_message"
	TimeParam      Param = "_time"
	LevelParam     Param = "_level"
	ComponentParam Param = "_component"
	TxIDParam      Param = "_txID"
)

type DriverInterface interface {
	log(data map[Param]string)
	beginTx(id TxID)
	endTx(id TxID)
}
