package logsystem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDriverManager_beginTx(t *testing.T) {
	m := NewManager()
	txID := m.beginTx(map[Param]string{})
	require.Equal(t, TxID(1), txID)
	txID = m.beginTx(map[Param]string{})
	require.Equal(t, TxID(2), txID)
}
