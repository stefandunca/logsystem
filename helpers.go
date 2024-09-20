package logsystem

import (
	"strconv"
	"strings"
	"time"
)

type KnownParams struct {
	Timestamp int64
	Level     string
	Message   string
	Component string
	TxID      string
}

func extractKnownParams(data map[Param]string) KnownParams {
	p := KnownParams{}
	if val, ok := data[TimeParam]; ok {
		parsedTime, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			p.Timestamp = parsedTime
		}
	}

	if p.Timestamp == 0 {
		p.Timestamp = time.Now().Unix()
	}

	if val, ok := data[LevelParam]; ok {
		p.Level = strings.ToUpper(val)
	}

	if val, ok := data[MessageParam]; ok {
		p.Message = val
	}

	if val, ok := data[ComponentParam]; ok {
		p.Component = val
	}

	if val, ok := data[TxIDParam]; ok {
		p.TxID = val
	}

	return p
}
