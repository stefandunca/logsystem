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
func (m *MockDriver) log(data map[Param]string) {
	m.Called(data)
}

func (m *MockDriver) beginTx(id TxID, attr map[Param]string) {
	m.Called(id, attr)
}

func (m *MockDriver) endTx(txID TxID) {
	m.Called(txID)
}

func (m *MockDriver) stop() {
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

	mockedDriverInterfaces := make([]DriverInterface, 0)
	for _, mockDriver := range mockDrivers {
		mockedDriverInterfaces = append(mockedDriverInterfaces, mockDriver)
	}
	manager.AddDrivers(mockedDriverInterfaces)

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

	for _, mockDriver := range mockDrivers {
		mockDriver.AssertExpectations(t)
	}
}
