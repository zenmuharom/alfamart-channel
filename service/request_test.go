package service

import (
	"alfamart-channel/util"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestConstructBit61(t *testing.T) {

	type TestCase struct {
		BillNumber  string
		ProductCode string
		Expected    string
	}

	err := util.LoadConfig("../.")

	logger := zenlogger.NewZenlogger()

	err = util.ConnectDB()
	if err != nil {
		logger.Error(err.Error())
		log.Fatalln(err)
		os.Exit(1)
	}

	testCases := []TestCase{
		{
			BillNumber:  "0211",
			ProductCode: "001001",
			Expected:    "0021000000001",
		},
		{
			BillNumber:  "0211234",
			ProductCode: "001001",
			Expected:    "0021000001234",
		},
		{
			BillNumber:  "00211234",
			ProductCode: "001001",
			Expected:    "0021000001234",
		},
		{
			BillNumber:  "1234",
			ProductCode: "030002",
			Expected:    "0000000001234",
		},
		{
			BillNumber:  "1001",
			ProductCode: "060100",
			Expected:    "0000000001234",
		},
	}

	alfamartCoreService := NewCoreService(logger)
	for _, tc := range testCases {

		result := alfamartCoreService.ConstructBit61(tc.ProductCode, tc.BillNumber)
		require.Equal(t, tc.Expected, result)
	}
}
