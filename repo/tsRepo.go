package repo

import (
	"alfamart-channel/util"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/zenmuharom/zenlogger"
)

type DefaultTSRepo struct {
	logger zenlogger.Zenlogger
}

type TSRepo interface {
	// Request(req dto.TsReq) (tsRes map[string]interface{}, err error)
	Request2(url string, timeout int, req map[string]interface{}) (tsRes map[string]interface{}, err error)
}

func NewTSRepo(logger zenlogger.Zenlogger) TSRepo {
	return &DefaultTSRepo{logger: logger}
}

func (repo *DefaultTSRepo) Request2(url string, timeout int, req map[string]interface{}) (tsRes map[string]interface{}, err error) {
	repo.logger.Debug("Switching", zenlogger.ZenField{Key: "url", Value: url}, zenlogger.ZenField{Key: "timeout", Value: timeout}, zenlogger.ZenField{Key: "req", Value: req})

	jsonData, _ := json.Marshal(req)
	repo.logger.Info("Switching request", zenlogger.ZenField{Key: "url", Value: url}, zenlogger.ZenField{Key: "body", Value: string(jsonData)})

	if url == "" {
		url = util.GetConfig().TS_URL
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		repo.logger.Error(err.Error())
		return tsRes, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Duration(timeout * int(time.Millisecond)),
	}
	response, err := client.Do(request)
	if err != nil {
		repo.logger.Error(err.Error())
		return tsRes, err
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		repo.logger.Info(err.Error())
		return tsRes, err
	}

	var responseBody map[string]interface{}

	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		repo.logger.Error(err.Error())
		return tsRes, err
	}
	repo.logger.Info("Switching response", zenlogger.ZenField{Key: "status", Value: response.Status}, zenlogger.ZenField{Key: "body", Value: responseBody})
	defer response.Body.Close()

	return responseBody, nil
}
