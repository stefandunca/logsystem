package logsystem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDriverManager_beginTx(t *testing.T) {
	m := NewManager(nil, Config{})
	txID := m.beginTx()
	require.Equal(t, TxID(1), txID)
	txID = m.beginTx()
	require.Equal(t, TxID(2), txID)
}
