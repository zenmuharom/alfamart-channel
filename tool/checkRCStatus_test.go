package tool

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestCheckRCStatus(t *testing.T) {

	logger := zenlogger.NewZenlogger()

	type TestCase struct {
		Rc       string
		Rcs      []byte
		Expected bool
	}

	testCases := []TestCase{
		{
			Rc:       "05",
			Rcs:      []byte(`["0", "00"]`),
			Expected: false,
		},
		{
			Rc:       "05",
			Rcs:      []byte(``),
			Expected: false,
		},
		{
			Rc:       "05",
			Rcs:      []byte(`[]`),
			Expected: false,
		},
		{
			Rc:       "",
			Rcs:      []byte(`[]`),
			Expected: false,
		},
		{
			Rc:       "68",
			Rcs:      []byte(`["68", "92"]`),
			Expected: true,
		},
		{
			Rc:       "68",
			Rcs:      []byte{},
			Expected: false,
		},
	}

	for _, tc := range testCases {

		result := CheckRCStatus(logger, tc.Rc, tc.Rcs)

		require.Equal(t, tc.Expected, result)

	}
}
