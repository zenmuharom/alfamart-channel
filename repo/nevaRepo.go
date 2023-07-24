package repo

import (
	"alfamart-channel/dto"
	"alfamart-channel/util"
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"net/http"

	"github.com/zenmuharom/zenlogger"
	"golang.org/x/net/html/charset"
)

type DefaultNevaRepo struct {
	logger zenlogger.Zenlogger
}

type NevaRepo interface {
	GetToken(req dto.NevaInputToken) (nevaRes dto.NevaGetTokenResponse, err error)
	DebetAccountV2(req dto.NevaInputTransaction) (nevaResponse dto.NevaDebetAccountV2Response, err error)
}

func NewNevaRepo(logger zenlogger.Zenlogger) NevaRepo {
	return &DefaultNevaRepo{logger: logger}
}

func (repo *DefaultNevaRepo) GetToken(req dto.NevaInputToken) (nevaResponse dto.NevaGetTokenResponse, err error) {

	var root = dto.NevaRoot{}
	root.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	root.Xsd = "http://www.w3.org/2001/XMLSchema"
	root.Env = "http://schemas.xmlsoap.org/soap/envelope/"
	root.Urn = "urn:maia"
	root.Header = dto.NevaHeader{}
	root.Body = dto.NevaBody{}

	root.Body.Request = dto.NevaGetTokenReq{
		EncodingStyle:  "http://schemas.xmlsoap.org/soap/encoding/",
		NevaInputToken: req,
	}

	out, _ := xml.MarshalIndent(&root, "", "")
	reqBody := string(out)
	repo.logger.Info("GetToken", zenlogger.ZenField{Key: "url", Value: util.GetConfig().NEVA_URL}, zenlogger.ZenField{Key: "request", Value: reqBody})

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	clientReq, err := http.NewRequest("POST", util.GetConfig().NEVA_URL, bytes.NewBufferString(reqBody))
	clientReq.Header.Set("Content-Type", "text/xml")
	if err != nil {
		repo.logger.Error("Fail to build request neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}
	response, err := client.Do(clientReq)
	if err != nil {
		repo.logger.Error("Fail to get response from neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&nevaResponse)
	if err != nil {
		repo.logger.Error("Fail to decode response from neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	nevaResponseDebug, _ := xml.Marshal(nevaResponse)
	repo.logger.Info("GetToken", zenlogger.ZenField{Key: "HTTPStatus", Value: response.StatusCode}, zenlogger.ZenField{Key: "body", Value: string(nevaResponseDebug)})
	defer response.Body.Close()
	if err != nil {
		repo.logger.Error("Fail to decode response from neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	return
}

func (repo *DefaultNevaRepo) DebetAccountV2(req dto.NevaInputTransaction) (nevaResponse dto.NevaDebetAccountV2Response, err error) {

	var root = dto.NevaRoot{}
	root.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	root.Xsd = "http://www.w3.org/2001/XMLSchema"
	root.Env = "http://schemas.xmlsoap.org/soap/envelope/"
	root.Urn = "urn:maia"
	root.Header = dto.NevaHeader{}
	root.Body = dto.NevaBody{}

	root.Body.Request = dto.NevaDebetAccountV2Req{
		EncodingStyle:        "http://schemas.xmlsoap.org/soap/encoding/",
		NevaInputTransaction: req,
	}

	out, _ := xml.MarshalIndent(&root, "", "")
	reqBody := string(out)
	repo.logger.Info("NEVA request", zenlogger.ZenField{Key: "url", Value: util.GetConfig().NEVA_URL}, zenlogger.ZenField{Key: "request", Value: reqBody})

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	clientReq, err := http.NewRequest("POST", util.GetConfig().NEVA_URL, bytes.NewBufferString(reqBody))
	clientReq.Header.Set("Content-Type", "text/xml")
	if err != nil {
		repo.logger.Error("Fail to build request neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}
	response, err := client.Do(clientReq)
	if err != nil {
		repo.logger.Error("Fail to get response from neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&nevaResponse)
	if err != nil {
		repo.logger.Error("Fail to decode response from neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	nevaResponseDebug, _ := xml.Marshal(nevaResponse)
	repo.logger.Info("NEVA response", zenlogger.ZenField{Key: "HTTPStatus", Value: response.StatusCode}, zenlogger.ZenField{Key: "body", Value: string(nevaResponseDebug)})
	defer response.Body.Close()
	if err != nil {
		repo.logger.Error("Fail to decode response from neva", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	return
}
