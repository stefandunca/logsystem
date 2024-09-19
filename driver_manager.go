package logsystem

import (
	"log"
	"sync/atomic"
)

type DriverManager struct {
	drivers []DriverInterface

	lastTxID atomic.Int64
}

func NewManager(factories []DriverFactoryInterface, config Config) *DriverManager {
	drivers := make([]DriverInterface, 0, len(factories))
	for _, factory := range factories {
		if _, ok := config.Drivers[factory.driverID()]; ok {
			driver := factory.createDriver(config.Drivers[factory.driverID()])
			if driver == nil {
				log.Printf("Failed to create driver %s", factory.driverID())
				continue
			}
			drivers = append(drivers, driver)
		}
	}

	return &DriverManager{
		drivers: drivers,
	}
}

func (m *DriverManager) log(data map[Param]string) {
	for _, driver := range m.drivers {
		driver.log(data)
	}
}

func (m *DriverManager) beginTx(attr map[Param]string) TxID {
	txID := TxID(m.lastTxID.Add(1))

	for _, driver := range m.drivers {
		driver.beginTx(txID, attr)
	}
	return txID
}

func (m *DriverManager) endTx(id TxID) {
	for _, driver := range m.drivers {
		driver.endTx(id)
	}
}
