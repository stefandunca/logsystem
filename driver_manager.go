package logsystem

import (
	"sync/atomic"
)

type DriverManager struct {
	drivers []DriverInterface

	lastTxID atomic.Int64
}

func NewManager() *DriverManager {
	return &DriverManager{}
}

func (m *DriverManager) AddDriver(driver DriverInterface) {
	m.drivers = append(m.drivers, driver)
}

func (m *DriverManager) AddDrivers(drivers []DriverInterface) {
	m.drivers = append(m.drivers, drivers...)
}

func (m *DriverManager) log(data map[Param]string) {
	for _, driver := range m.drivers {
		driver.Log(data)
	}
}

func (m *DriverManager) beginTx(attr map[Param]string) TxID {
	txID := TxID(m.lastTxID.Add(1))

	for _, driver := range m.drivers {
		driver.BeginTx(txID, attr)
	}
	return txID
}

func (m *DriverManager) endTx(id TxID) {
	for _, driver := range m.drivers {
		driver.EndTx(id)
	}
}

func (m *DriverManager) stop() {
	for _, driver := range m.drivers {
		driver.Stop()
	}
}
