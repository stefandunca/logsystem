package logsystem

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDriver struct {
	mock.Mock
}

// MockDriver mocks DriverInterface
func (m *MockDriver) Log(data map[Param]string) {
	m.Called(data)
}

func (m *MockDriver) BeginTx(id TxID, attr map[Param]string) {
	m.Called(id, attr)
}

func (m *MockDriver) EndTx(txID TxID) {
	m.Called(txID)
}

func (m *MockDriver) Stop() {
	m.Called()
}

// TestDriverManagerIntegration tests the integration of the DriverManager with the DriverInterface
func TestDriverManagerIntegration(t *testing.T) {
	mockDrivers := make([]*MockDriver, 0)
	for i := 0; i < 2; i++ {
		mockDriver := &MockDriver{}

		mockDrivers = append(mockDrivers, mockDriver)
	}

	// Call manager API and validate expectations
	//

	manager := NewManager()
	require.Equal(t, 0, len(manager.drivers))

	for _, mockDriver := range mockDrivers {
		manager.AddDriver(mockDriver)
	}

	simpleLogPayload := map[Param]string{
		TimeParam:      strconv.FormatInt(time.Now().Unix(), 10),
		MessageParam:   "message",
		LevelParam:     "info",
		ComponentParam: "component",
	}

	for _, mockDriver := range mockDrivers {
		mockDriver.On("Log", simpleLogPayload).Once()
		mockDriver.On("BeginTx", mock.AnythingOfType("logsystem.TxID"), mock.AnythingOfType("map[logsystem.Param]string")).Once()
	}
	manager.log(simpleLogPayload)
	txID := manager.beginTx(map[Param]string{})

	txLogPayload := map[Param]string{
		TxIDParam: txID.String(),
	}
	for _, mockDriver := range mockDrivers {
		mockDriver.On("Log", txLogPayload).Once()
		mockDriver.On("EndTx", txID).Once()
		mockDriver.On("Stop").Once()
	}
	manager.log(txLogPayload)
	manager.endTx(txID)

	manager.stop()

	// Validate expectations for all mocks
	//

	for _, mockDriver := range mockDrivers {
		mockDriver.AssertExpectations(t)
	}
}
