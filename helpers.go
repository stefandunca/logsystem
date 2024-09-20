package logsystem

import (
	"fmt"
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

func formatLine(data map[Param]string, userFriendly bool) string {
	p := extractKnownParams(data)

	formattedTime := ""
	if userFriendly {
		t := time.Unix(p.Timestamp, 0)
		formattedTime = t.Format("[2006-01-02 15:04:05] ")
	} else {
		formattedTime = fmt.Sprintf("[%-10d] ", p.Timestamp)
	}

	optional := ""
	if p.TxID != "" || p.Component != "" {
		if p.Component != "" {
			optional = fmt.Sprintf("; Comp=[%s]", p.Component)
		}
		if p.TxID != "" {
			optional = fmt.Sprintf("%s; TxID=[%s]", optional, p.TxID)
		}
	}

	return fmt.Sprintf("%s%-5s %s%s", formattedTime, p.Level, p.Message, optional)
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
