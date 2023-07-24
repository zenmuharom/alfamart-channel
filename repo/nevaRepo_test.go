package repo

import (
	"alfamart-channel/dto"
	"alfamart-channel/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestGetToken(t *testing.T) {

	err := util.LoadConfig("../.")

	logger := zenlogger.NewZenlogger()
	nevaRepo := NewNevaRepo(logger)
	req := dto.NevaInputToken{
		Type: "urn:getToken",
		Username: dto.NevaField{
			Type:  "xsd:string",
			Value: "3m0N",
		},
		Password: dto.NevaField{
			Type:  "xsd:string",
			Value: "emonP@ssWd",
		},
		MitraCo: dto.NevaField{
			Type:  "xsd:string",
			Value: "148",
		},
	}
	res, err := nevaRepo.GetToken(req)
	require.NoError(t, err, "Gagal")
	require.Equal(t, "0000", res.Body.GetTokenResponse.OutputToken.Status.ResultCode)
	fmt.Println(fmt.Sprintf("%#v", res))

}

func TestDebetAccountV2(t *testing.T) {
	err := util.LoadConfig("../.")

	logger := zenlogger.NewZenlogger()
	nevaRepo := NewNevaRepo(logger)
	req := dto.NevaInputTransaction{
		Type: "urn:inputTransaction",
		Description: dto.NevaField{
			Type:  "xsd:string",
			Value: "beli",
		},
		Dest1Acc: dto.NevaField{
			Type:  "xsd:string",
			Value: "+628177022101910210485664",
		},
		Dest1Amount: dto.NevaField{
			Type:  "xsd:string",
			Value: "16500",
		},
		Phoneno: dto.NevaField{
			Type:  "xsd:string",
			Value: "+6281288313488",
		},
		Sessionid: dto.NevaField{
			Type:  "xsd:string",
			Value: "2023031310321955CD5F403992FF9D4D966E33386A6A19A158B1F0220C018256FF4649D3D01B82",
		},
		Source1Acc: dto.NevaField{
			Type:  "xsd:string",
			Value: "+6281288313488",
		},
		Source1Amount: dto.NevaField{
			Type:  "xsd:string",
			Value: "16500",
		},
		TransactionType: dto.NevaField{
			Type:  "xsd:string",
			Value: "Payment",
		},
		Mitraco: dto.NevaField{
			Type:  "xsd:string",
			Value: "148",
		},
		Prodcode: dto.NevaField{
			Type:  "xsd:string",
			Value: "001001",
		},
		Billingno: dto.NevaField{
			Type:  "xsd:string",
			Value: "0021008455002",
		},
		TraxId: dto.NevaField{
			Type:  "xsd:string",
			Value: "1678179415513",
		},
	}
	res, err := nevaRepo.DebetAccountV2(req)
	require.NoError(t, err, "Gagal")
	require.Equal(t, "0", res.Body.DebetAccountV2Response.OutputTransaction.Status.ResultCode)
	fmt.Println(fmt.Sprintf("%#v", res))
}
