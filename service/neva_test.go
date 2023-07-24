package service

import (
	"alfamart-channel/dto"
	"alfamart-channel/util"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestDebet(t *testing.T) {
	err := util.LoadConfig("../.")

	logger := zenlogger.NewZenlogger()

	nevaService := NewNevaService(logger)
	debetReq := dto.NevaDebet{
		Username:  "3m0N",
		Password:  "emonP@ssWd",
		MitraCo:   "148",
		SourceAcc: "+6281288313488",
		DestAcc:   "+628177022101910210485664",
		Amount:    "16500",
		PhoneNo:   "+6281288313488",
		ProdCode:  "001001",
		BillingNo: "0021008455002",
		TraxId:    "1678179415513",
	}
	status, err := nevaService.Debet(debetReq)
	require.NoError(t, err, "error")
	require.Equal(t, true, status)

}
