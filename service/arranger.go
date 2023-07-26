package service

import (
	"alfamart-channel/repo"

	"github.com/zenmuharom/zenlogger"
)

type DefaultArranger struct {
	logger zenlogger.Zenlogger
}

type ArrangerService interface {
	Arrange(productCode string, endpoint string, toArrange map[string]interface{}) (arranged []interface{}, err error)
}

func NewArrangerService(logger zenlogger.Zenlogger) ArrangerService {
	return &DefaultArranger{logger: logger}
}

func (service *DefaultArranger) Arrange(productCode string, endpoint string, toArrange map[string]interface{}) (arranged []interface{}, err error) {

	arranged = make([]interface{}, 0)

	serverResponseRepo := repo.NewServerResponseRepo(service.logger)
	srConfs, err := serverResponseRepo.FindAllByProductCodeAndEndpointParentOnly(productCode, endpoint)
	if err != nil {
		service.logger.Error("Arrange", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	for _, config := range srConfs {
		arranged = append(arranged, toArrange[config.Field.String])
	}

	return
}
