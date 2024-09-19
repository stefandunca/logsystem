package logsystem

import (
	"encoding/json"
)

type Config struct {
	defaultParams map[Param]string             `json:"defaultParams"`
	drivers       map[DriverID]json.RawMessage `json:"drivers"`
}
