package repo

import (
	"alfamart-channel/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestFeeFind(t *testing.T) {
	err := util.LoadConfig("../.")
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	err = util.ConnectDB()
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	logger := zenlogger.NewZenlogger()
	feeRepo := NewFeeRepo(logger)

	fee, err := feeRepo.Find("finnetDev", "001001", "20000")
	if err != nil {
		require.NoError(t, err, err.Error())
	}
	fmt.Println(fmt.Sprintf("%v", fee))

}
