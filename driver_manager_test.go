package logsystem

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type FailingDriverFactory struct {
}

func (f *FailingDriverFactory) driverID() DriverID {
	return "failing"
}

func (f *FailingDriverFactory) createDriver(config json.RawMessage) DriverInterface {
	return nil
}

func TestDriverManager_FailToCreateDriver(t *testing.T) {
	m := NewManager([]DriverFactoryInterface{
		&FailingDriverFactory{},
	}, Config{})
	require.Len(t, m.drivers, 0)
}

func TestDriverManager_beginTx(t *testing.T) {
	m := NewManager(nil, Config{})
	txID := m.beginTx()
	require.Equal(t, TxID(1), txID)
	txID = m.beginTx()
	require.Equal(t, TxID(2), txID)
}
