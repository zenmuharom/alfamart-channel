package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPadLeft(t *testing.T) {
	res := PadLeft("0021082999003", 13, "0")
	require.Equal(t, "0021082999003", res)
}

func TestPadRight(t *testing.T) {
	res := PadRight("0021082999003", 14, "0")
	require.Equal(t, "00210829990030", res)
}
