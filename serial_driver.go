package logsystem

import (
	"encoding/json"
	"sync"
)

const SerialDriverIDPostfix = "-serial"

// SerialDriverFactory implements DriverFactoryInterface
type SerialDriverFactory struct {
	provider DriverFactoryInterface
}

func NewSerialDriverFactory(provider DriverFactoryInterface) *SerialDriverFactory {
	_, ok := provider.(*SerialDriverFactory)
	if ok {
		panic("SerialDriverFactory cannot be nested")
	}

	return &SerialDriverFactory{
		provider: provider,
	}
}

func (f *SerialDriverFactory) DriverID() DriverID {
	return DriverID(string(f.provider.DriverID()) + SerialDriverIDPostfix)
}

func (f *SerialDriverFactory) CreateDriver(config json.RawMessage) (DriverInterface, error) {
	return f.provider.CreateDriver(config)
}

// SerialDriver implements DriverInterface
type SerialDriver struct {
	provider DriverInterface
	mutex    sync.Mutex
}

func NewSerialDriver(provider DriverInterface) *SerialDriver {
	_, ok := provider.(*SerialDriver)
	if ok {
		panic("SerialDriver cannot be nested")
	}

	return &SerialDriver{
		provider: provider,
	}
}

func (d *SerialDriver) Log(data map[Param]string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.provider.Log(data)
}

func (d *SerialDriver) BeginTx(id TxID, attr map[Param]string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.provider.BeginTx(id, attr)
}

func (d *SerialDriver) EndTx(id TxID) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.provider.EndTx(id)
}

func (d *SerialDriver) Stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.provider.Stop()
}
