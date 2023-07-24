package service

import (
	"alfamart-channel/util"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestIdentifySignature(t *testing.T) {
	type TestCase struct {
		Req   map[string]interface{}
		Valid bool
	}

	logger := zenlogger.NewZenlogger()

	err := util.LoadConfig("../.")

	require.NoError(t, err, "error occured when load config")

	err = util.ConnectDB()
	if err != nil {
		logger.Error(err.Error())
		log.Fatalln(err)
		os.Exit(1)
	}

	require.NoError(t, err, "error occured when connecting to DB")

	testCases := []TestCase{
		{
			Req: map[string]interface{}{
				"userName":         "testTS",
				"signature":        "",
				"productCode":      "080015",
				"terminal":         "G009TRKK",
				"terminalName":     "FINNET",
				"terminalLocation": "Bidakara Pancoran",
				"transactionType":  "38",
				"channel":          "6018",
				"billNumber":       "7302046706800014",
				"amount":           "201600",
				"feeAmount":        "0",
				"bit32":            "770",
				"bit33":            "770924",
				"bit61":            "7302046706800014",
				"traxId":           "1689056039763",
				"timeStamp":        "11-07-2023 13:13:59:+07:00",
			},
			Valid: false,
		},
		{
			Req: map[string]interface{}{
				"userName":         "testTS",
				"signature":        true,
				"productCode":      "080015",
				"terminal":         "G009TRKK",
				"terminalName":     "FINNET",
				"terminalLocation": "Bidakara Pancoran",
				"transactionType":  "38",
				"channel":          "6018",
				"billNumber":       "7302046706800014",
				"amount":           "201600",
				"feeAmount":        "0",
				"bit32":            "770",
				"bit33":            "770924",
				"bit61":            "7302046706800014",
				"traxId":           "1689056039763",
				"timeStamp":        "11-07-2023 13:13:59:+07:00",
			},
			Valid: false,
		},
		{
			Req: map[string]interface{}{
				"userName":         "testTS",
				"signature":        1,
				"productCode":      "080015",
				"terminal":         "G009TRKK",
				"terminalName":     "FINNET",
				"terminalLocation": "Bidakara Pancoran",
				"transactionType":  "38",
				"channel":          "6018",
				"billNumber":       "7302046706800014",
				"amount":           "201600",
				"feeAmount":        "0",
				"bit32":            "770",
				"bit33":            "770924",
				"bit61":            "7302046706800014",
				"traxId":           "1689056039763",
				"timeStamp":        "11-07-2023 13:13:59:+07:00",
			},
			Valid: false,
		},
		{
			Req: map[string]interface{}{
				"userName":         "testTS",
				"signature":        "55072f08c73ef01c29f2ff52923cf9ef",
				"productCode":      "080015",
				"terminal":         "G009TRKK",
				"terminalName":     "FINNET",
				"terminalLocation": "Bidakara Pancoran",
				"transactionType":  "38",
				"channel":          "6018",
				"billNumber":       "7302046706800014",
				"amount":           "201600",
				"feeAmount":        "0",
				"bit32":            "770",
				"bit33":            "770924",
				"bit61":            "7302046706800014",
				"traxId":           "1689057690767",
				"timeStamp":        "11-07-2023 13:41:30:+07:00",
			},
			Valid: true,
		},
		{
			Req: map[string]interface{}{
				"userName":         "testTS",
				"signature":        "8c26ff34940644940543b8de4b3fb420ce3c66afc8c8ad2d66254e6f4e92c93e",
				"productCode":      "080015",
				"terminal":         "G009TRKK",
				"terminalName":     "FINNET",
				"terminalLocation": "Bidakara Pancoran",
				"transactionType":  "38",
				"channel":          "6018",
				"billNumber":       "7302046706800014",
				"amount":           "201600",
				"feeAmount":        "0",
				"bit32":            "770",
				"bit33":            "770924",
				"bit61":            "7302046706800014",
				"traxId":           "1689057726592",
				"timeStamp":        "11-07-2023 13:42:06:+07:00",
			},
			Valid: true,
		},
		{
			Req: map[string]interface{}{
				"userName":         "testTS",
				"signature":        "INVALID8c73ef01c29f2ff52923cf9ef",
				"productCode":      "080015",
				"terminal":         "G009TRKK",
				"terminalName":     "FINNET",
				"terminalLocation": "Bidakara Pancoran",
				"transactionType":  "38",
				"channel":          "6018",
				"billNumber":       "7302046706800014",
				"amount":           "201600",
				"feeAmount":        "0",
				"bit32":            "770",
				"bit33":            "770924",
				"bit61":            "7302046706800014",
				"traxId":           "1689057690767",
				"timeStamp":        "11-07-2023 13:41:30:+07:00",
			},
			Valid: false,
		},
		{
			Req: map[string]interface{}{
				"userName":         "testTS",
				"signature":        "INVALID4940644940543b8de4b3fb420ce3c66afc8c8ad2d66254e6f4e92c93e",
				"productCode":      "080015",
				"terminal":         "G009TRKK",
				"terminalName":     "FINNET",
				"terminalLocation": "Bidakara Pancoran",
				"transactionType":  "38",
				"channel":          "6018",
				"billNumber":       "7302046706800014",
				"amount":           "201600",
				"feeAmount":        "0",
				"bit32":            "770",
				"bit33":            "770924",
				"bit61":            "7302046706800014",
				"traxId":           "1689057726592",
				"timeStamp":        "11-07-2023 13:42:06:+07:00",
			},
			Valid: false,
		},
	}

	signatureService := NewSignatureService(logger)

	for _, tc := range testCases {
		valid, err := signatureService.Check("/", tc.Req)
		require.NoError(t, err, fmt.Sprintf("error occured when checking: %v", tc.Req["signature"]))
		require.Equal(t, tc.Valid, valid, fmt.Sprintf("Expected:%v got:%v when test %v", tc.Valid, valid, tc.Req["signature"]))
	}
}
