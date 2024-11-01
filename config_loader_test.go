package logsystem

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type FailingDriverFactory struct {
}

func (f *FailingDriverFactory) DriverID() DriverID {
	return "failing"
}

func (f *FailingDriverFactory) CreateDriver(config json.RawMessage) (DriverInterface, error) {
	return nil, errors.New("failed to create driver")
}

type SuccessDriverFactory struct {
	failCounts int
}

func (f *SuccessDriverFactory) DriverID() DriverID {
	return "success"
}

func (f *SuccessDriverFactory) CreateDriver(config json.RawMessage) (DriverInterface, error) {
	return &SuccessDriverFactory{}, nil
}

func (f *SuccessDriverFactory) beginTx(id TxID, attr map[Param]string) {
	panic("Not expected to be called")
}

func (f *SuccessDriverFactory) endTx(txID TxID) {
	panic("Not expected to be called")
}

func (f *SuccessDriverFactory) log(map[Param]string) {
	f.failCounts += 1
}

func (f *SuccessDriverFactory) stop() {
	panic("Not expected to be called")
}

func TestConfigLoader_FailToCreateDriver(t *testing.T) {
	m, err := CreateLogManagerWithConfig(
		[]DriverFactoryInterface{
			&FailingDriverFactory{},
		},
		Config{
			Drivers: map[DriverID]json.RawMessage{
				"failing": json.RawMessage(`{}`),
			},
		},
	)
	require.ErrorIs(t, err, ErrorAllDriversFailed)
	require.Nil(t, m)
}

func TestConfigLoader_PartiallyFailToCreateDriver(t *testing.T) {
	m, err := CreateLogManagerWithConfig(
		[]DriverFactoryInterface{
			&FailingDriverFactory{},
			&SuccessDriverFactory{},
		},
		Config{
			Drivers: map[DriverID]json.RawMessage{
				"failing": json.RawMessage(`{}`),
				"success": json.RawMessage(`{}`),
			},
		},
	)
	require.ErrorIs(t, err, ErrorSomeDriversFailed)
	require.Len(t, m.drivers, 1)
	require.Equal(t, 1, m.drivers[0].(*SuccessDriverFactory).failCounts)
}

func TestConfigLoader_SuccessfullyCreateDriver(t *testing.T) {
	m, err := CreateLogManagerWithConfig(
		[]DriverFactoryInterface{
			&SuccessDriverFactory{},
		},
		Config{
			Drivers: map[DriverID]json.RawMessage{
				"success": json.RawMessage(`{}`),
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, m.drivers, 1)
}

// MockDriverFactory mocks DriverFactoryInterface
type MockDriverFactory struct {
	mock.Mock
	mockDriver *MockDriver
}

func NewMockDriverFactory(mockDriver *MockDriver) *MockDriverFactory {
	return &MockDriverFactory{mockDriver: mockDriver}
}

func (m *MockDriverFactory) DriverID() DriverID {
	args := m.Called()
	return args.Get(0).(DriverID)
}

func (m *MockDriverFactory) CreateDriver(config json.RawMessage) (DriverInterface, error) {
	m.Called(config)
	return m.mockDriver, nil
}

// TestConfigLoaderIntegration tests the integration of the DriverManager with the DriverFactoryInterface
func TestConfigLoaderIntegration(t *testing.T) {
	// Prepare dummy driver and factory
	//

	dummyConfig := json.RawMessage{}
	err := dummyConfig.UnmarshalJSON([]byte(`{"key":"value"}`))
	require.NoError(t, err)

	mockDrivers := make([]*MockDriver, 0)
	mockFactories := make([]DriverFactoryInterface, 0)
	for i := 0; i < 2; i++ {
		mockDriver := &MockDriver{}

		mockDrivers = append(mockDrivers, mockDriver)

		mockID := DriverID("mockID" + strconv.Itoa(i))
		mockFactory := NewMockDriverFactory(mockDriver)
		mockFactory.On("DriverID").Return(mockID)
		mockFactory.On("CreateDriver", dummyConfig).Return(mockDriver)

		mockFactories = append(mockFactories, mockFactory)
	}

	// Call manager API and validate expectations
	//

	manager, err := CreateLogManagerWithConfig(mockFactories, Config{
		Drivers: map[DriverID]json.RawMessage{
			mockFactories[0].DriverID(): dummyConfig,
			mockFactories[1].DriverID(): dummyConfig,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(manager.drivers))
	require.Equal(t, 2, len(mockDrivers))

	simpleLogPayload := map[Param]string{
		TimeParam:      strconv.FormatInt(time.Now().Unix(), 10),
		MessageParam:   "message",
		LevelParam:     "info",
		ComponentParam: "component",
	}

	for _, mockDriver := range mockDrivers {
		mockDriver.On("log", simpleLogPayload).Once()
		mockDriver.On("beginTx", mock.AnythingOfType("logsystem.TxID"), mock.AnythingOfType("map[logsystem.Param]string")).Once()
	}
	manager.log(simpleLogPayload)
	txID := manager.beginTx(map[Param]string{})

	txLogPayload := map[Param]string{
		TxIDParam: txID.String(),
	}
	for _, mockDriver := range mockDrivers {
		mockDriver.On("log", txLogPayload).Once()
		mockDriver.On("endTx", txID).Once()
		mockDriver.On("stop").Once()
	}
	manager.log(txLogPayload)
	manager.endTx(txID)

	manager.stop()

	// Validate expectations for all mocks
	//

	for _, mockFactoryIf := range mockFactories {
		mockFactory := mockFactoryIf.(*MockDriverFactory)
		mockFactory.AssertExpectations(t)
	}
	for _, mockDriver := range mockDrivers {
		mockDriver.AssertExpectations(t)
	}
}
