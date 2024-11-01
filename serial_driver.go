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

func (d *SerialDriver) log(data map[Param]string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.provider.log(data)
}

func (d *SerialDriver) beginTx(id TxID, attr map[Param]string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.provider.beginTx(id, attr)
}

func (d *SerialDriver) endTx(id TxID) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.provider.endTx(id)
}
