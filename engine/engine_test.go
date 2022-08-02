package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	_, err := NewEngine("mysql", "root", "hcwnbs", "localhost:3306", "ut-test")
	require.Nil(t, err)
}