package repo

import (
	"alfamart-channel/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestRouteFind(t *testing.T) {
	err := util.LoadConfig("../.")
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	err = util.ConnectDB()
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	logger := zenlogger.NewZenlogger()
	routeRepo := NewRouteRepo(logger)

	route, err := routeRepo.FindAll()
	if err != nil {
		require.NoError(t, err, err.Error())
	}
	fmt.Println(fmt.Sprintf("%v", route))

}
