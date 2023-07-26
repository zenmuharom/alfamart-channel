package repo

import (
	"alfamart-channel/util"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestServerRequest_FindAllByEndpoint(t *testing.T) {
	err := util.LoadConfig("../.")
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	err = util.ConnectDB()
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	logger := zenlogger.NewZenlogger()
	serverRequestRepo := NewServerRequestRepo(logger)

	serverRequestConfs, err := serverRequestRepo.FindAllByEndpoint("inquiryData")
	require.NoError(t, err, "error when call FindAllByEndpoint")

	for _, srC := range serverRequestConfs {
		require.NotEmpty(t, srC, "dont empty")
	}
}

func TestServerRequest_FindByEndpointAndFieldAs(t *testing.T) {
	err := util.LoadConfig("../.")
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	err = util.ConnectDB()
	if err != nil {
		require.NoError(t, err, err.Error())
	}

	logger := zenlogger.NewZenlogger()
	serverRequestRepo := NewServerRequestRepo(logger)

	signature, err := serverRequestRepo.FindByEndpointAndFieldAs("inquiryData", "signature")
	require.NoError(t, err, "error when call FindByEndpointAndFieldAs")
	require.Equal(t, "Signature", signature.Field.String)
	require.True(t, signature.Required.Bool, "not same")

}
