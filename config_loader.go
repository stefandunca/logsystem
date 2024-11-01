package logsystem

import "errors"

type failedDriver struct {
	id  DriverID
	err error
}

var (
	ErrorAllDriversFailed  = errors.New("all drivers failed to initialize")
	ErrorSomeDriversFailed = errors.New("some drivers failed to initialize")
)

func CreateLogManagerWithConfig(factories []DriverFactoryInterface, config Config) (*DriverManager, error) {
	drivers, failedDrivers := matchConfigWithDrivers(factories, config)
	if (len(drivers) == 0) && (len(failedDrivers) > 0) {
		return nil, ErrorAllDriversFailed
	}
	mgr := NewManager()
	mgr.AddDrivers(drivers)
	if len(failedDrivers) > 0 {
		for _, failedDriver := range failedDrivers {
			mgr.log(map[Param]string{
				"driver_id": string(failedDriver.id),
				"error":     failedDriver.err.Error(),
			})
		}
		return mgr, ErrorSomeDriversFailed
	}
	return mgr, nil
}

func matchConfigWithDrivers(factories []DriverFactoryInterface, config Config) (drivers []DriverInterface, failedDrivers []failedDriver) {
	failedDrivers = make([]failedDriver, 0)
	drivers = make([]DriverInterface, 0)

	for _, factory := range factories {
		if _, ok := config.Drivers[factory.DriverID()]; ok {
			driver, crErr := factory.CreateDriver(config.Drivers[factory.DriverID()])
			if crErr != nil {
				failedDrivers = append(failedDrivers, failedDriver{
					id:  factory.DriverID(),
					err: crErr,
				})
				continue
			}
			drivers = append(drivers, driver)
		}
	}

	return drivers, failedDrivers
}
