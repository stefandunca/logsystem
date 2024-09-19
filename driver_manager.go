package logsystem

import "sync/atomic"

type DriverManager struct {
	drivers []DriverInterface

	lastTxID atomic.Int64
}

func NewManager(factories []DriverFactoryInterface, config Config) *DriverManager {
	drivers := make([]DriverInterface, 0, len(factories))
	for _, factory := range factories {
		if _, ok := config.Drivers[factory.driverID()]; ok {
			driver := factory.createDriver(config.Drivers[factory.driverID()])
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

func (m *DriverManager) beginTx() TxID {
	txID := TxID(m.lastTxID.Add(1))

	for _, driver := range m.drivers {
		driver.beginTx(txID)
	}
	return txID
}

func (m *DriverManager) endTx(id TxID) {
	for _, driver := range m.drivers {
		driver.endTx(id)
	}
}
