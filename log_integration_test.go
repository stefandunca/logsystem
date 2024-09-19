package logsystem

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockDriverFactory mocks DriverFactoryInterface
type MockDriverFactory struct {
	mock.Mock
	mockDriver *MockDriver
}

func NewMockDriverFactory(mockDriver *MockDriver) *MockDriverFactory {
	return &MockDriverFactory{mockDriver: mockDriver}
}

func (m *MockDriverFactory) driverID() DriverID {
	args := m.Called()
	return args.Get(0).(DriverID)
}

func (m *MockDriverFactory) createDriver(config json.RawMessage) DriverInterface {
	m.Called(config)
	return m.mockDriver
}

type MockDriver struct {
	mock.Mock
}

// MockDriver mocks DriverInterface
func (m *MockDriver) log(data map[Param]string) {
	m.Called(data)
}

func (m *MockDriver) beginTx(id TxID) {
	m.Called(id)
}

func (m *MockDriver) endTx(txID TxID) {
	m.Called(txID)
}

// TestDriverManagerIntegration tests the integration of the DriverManager with the DriverInterface
func TestDriverManagerIntegration(t *testing.T) {
	//
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
		mockFactory.On("driverID").Return(mockID)
		mockFactory.On("createDriver", dummyConfig).Return(mockDriver)

		mockFactories = append(mockFactories, mockFactory)
	}

	//
	// Validate assumptions
	//

	manager := NewManager(mockFactories, Config{
		drivers: map[DriverID]json.RawMessage{
			mockFactories[0].driverID(): dummyConfig,
			mockFactories[1].driverID(): dummyConfig,
		},
	})
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
		mockDriver.On("beginTx", mock.AnythingOfType("logsystem.TxID")).Once()
	}
	manager.log(simpleLogPayload)
	txID := manager.beginTx()

	txLogPayload := map[Param]string{
		TxIDParam: txID.String(),
	}
	for _, mockDriver := range mockDrivers {
		mockDriver.On("log", txLogPayload).Once()
		mockDriver.On("endTx", txID).Once()
	}
	manager.log(txLogPayload)
	manager.endTx(txID)

	for _, mockFactoryIf := range mockFactories {
		mockFactory := mockFactoryIf.(*MockDriverFactory)
		mockFactory.AssertExpectations(t)
	}
	for _, mockDriver := range mockDrivers {
		mockDriver.AssertExpectations(t)
	}
}