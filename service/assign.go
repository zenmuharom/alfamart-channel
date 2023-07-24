package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/function"
	"alfamart-channel/repo"
	"alfamart-channel/tool"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/zenmuharom/zenlogger"
)

type DefaultAssignService struct {
	logger zenlogger.Zenlogger
}

type AssignService interface {
	AssignServerResponse(productCode, endpoint string, middlewResponseIDs []int, middlwareResponse map[string]interface{}) (serverResponse map[string]interface{}, err error)
	AssignMiddlewareRequest(middleware, process, productCode string, serverRequest map[string]interface{}) (middlewareRequest map[string]interface{}, err error)
}

func NewAssignService(logger zenlogger.Zenlogger) AssignService {
	return &DefaultAssignService{logger: logger}
}

func (service *DefaultAssignService) AssignMiddlewareRequest(middleware, process, productCode string, serverRequest map[string]interface{}) (middlewareRequest map[string]interface{}, err error) {
	service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "middleware", Value: middleware}, zenlogger.ZenField{Key: "serverRequest", Value: serverRequest})

	// Find all middleware request
	middlewareRequestRepo := repo.NewMiddlewareRequestRepo(service.logger)
	middlewareRequestConfs, err := middlewareRequestRepo.FindAllByAttributes(domain.MiddlewareRequest{
		Middleware:  sql.NullString{String: middleware},
		Process:     sql.NullString{String: process},
		ProductCode: sql.NullString{String: productCode},
	})
	service.logger.Debug("middlewareRequestConfs", zenlogger.ZenField{Key: "data", Value: middlewareRequestConfs})

	// gather all server response fieldIDs
	fieldIDs := make([]int, 0)
	for _, middlewareRequest := range middlewareRequestConfs {
		fieldIDs = append(fieldIDs, middlewareRequest.Id)
	}

	// Find all middleware request value config
	mrvRepo := repo.NewMiddlewareRequestValueRepo(service.logger)
	// mrvConfs, err := mrvRepo.FindAllByAttributes(domain.MiddlewareRequestValue{CaCode: sql.NullString{String: caCode}})
	mrvConfs, err := mrvRepo.FindAllByFieldIDs(fieldIDs)
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	assigner := function.NewAssigner2(service.logger)
	requestValues, errs := assigner.MiddlewareRequestParse(mrvConfs, serverRequest)
	if len(errs) > 0 {
		err = errors.New("Occur some errors")
		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
		return
	} else {
		service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "requestValues", Value: requestValues})
	}

	// FOR DEBUG in local
	// fmt.Println(service.logger.Debug("sebelum diconsttruct", zenlogger.ZenField{Key: "requestValues", Value: requestValues}))

	middlewareRequest, errs = assigner.MiddlewareRequestConstruct(middlewareRequestConfs, requestValues)
	if len(errs) > 0 {
		err = errors.New("Occur some errors")
		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
		return
	} else {
		service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "middlewareRequest", Value: middlewareRequest})
	}

	service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "middlewareRequest", Value: middlewareRequest}, zenlogger.ZenField{Key: "requestValues", Value: requestValues})

	return
}

func (service *DefaultAssignService) AssignServerResponse(productCode, endpoint string, middlewResponseIDs []int, middlewareResponse map[string]interface{}) (serverResponse map[string]interface{}, err error) {
	service.logger.Debug("AssignServerResponse", zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "middlewareResponse", Value: middlewareResponse})

	var serverResponseConfs []domain.ServerResponse

	// find all server response
	serverResponseRepo := repo.NewServerResponseRepo(service.logger)
	// if productCode is not empty select it base on productCode otherwise select the default one
	if productCode != "" {
		serverResponseConfs, err = serverResponseRepo.FindAllByProductCodeAndEndpoint(productCode, endpoint)
		// if config is zero select config from default
		if len(serverResponseConfs) == 0 {
			serverResponseConfs, err = serverResponseRepo.FindAllByProductCodeAndEndpoint(productCode, endpoint)
		}
	} else {
		serverResponseConfs, err = serverResponseRepo.FindAllByProductCodeAndEndpoint(productCode, endpoint)
	}
	service.logger.Debug("serverResponseConfs", zenlogger.ZenField{Key: "data", Value: serverResponseConfs})
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	// gather all server response fieldIDs
	fieldIDs := make([]int, 0)
	for _, serverResponse := range serverResponseConfs {
		fieldIDs = append(fieldIDs, serverResponse.Id)
	}

	// Find all server response value config
	srvRepo := repo.NewServerResponseValueRepo(service.logger)
	// srvConfs, err := srvRepo.FindAllByProductCode(productCode)
	srvConfs, err := srvRepo.FindAllByServerResponseID(productCode, fieldIDs, middlewResponseIDs)
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	assigner := function.NewAssigner2(service.logger)

	responseValues, errs := assigner.ServerResponseParse(srvConfs, middlewareResponse)
	if len(errs) > 0 {
		err = errors.New("Occur some errors")
		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
		return
	}

	serverResponse, errs = assigner.ServerResponseConstruct(serverResponseConfs, responseValues)
	if len(errs) > 0 {
		err = errors.New("Occur some errors")
		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
		return
	}

	// find rc key
	resultCodeKey, resultCodeIValue := tool.FindFieldAs(service.logger, domain.SERVER_RESPONSE, endpoint, "resultCode", serverResponse)
	rcConfigRepo := repo.NewRcConfigRepo(service.logger)
	rcCode, _ := strconv.ParseInt(fmt.Sprintf("%v", resultCodeIValue), 10, 64)
	rcConfig, err := rcConfigRepo.FindRc(productCode, rcCode)

	if resultCodeKey.Parent == "" {
		serverResponse[resultCodeKey.Field] = rcConfig.Code.String
	} else {
		// if it has parent key
		parentRCObject := serverResponse[resultCodeKey.Parent]
		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
		newParentResultCodeObject := make(map[string]interface{})
		switch valueOfVariableParentRC.Kind() {
		case reflect.Map:
			iter := valueOfVariableParentRC.MapRange()
			for iter.Next() {
				if iter.Key().String() == resultCodeKey.Field {
					newParentResultCodeObject[iter.Key().String()] = rcConfig.Code.String
				} else {
					newParentResultCodeObject[iter.Key().String()] = iter.Value().Interface()
				}
			}
		}
		serverResponse[resultCodeKey.Parent] = newParentResultCodeObject
	}

	resultDescKey, _ := tool.FindFieldAs(service.logger, domain.SERVER_RESPONSE, endpoint, "resultDesc", serverResponse)
	if resultDescKey.Parent == "" {
		serverResponse[resultDescKey.Field] = rcConfig.DescEng.String
	} else {
		// if it has parent key
		parentRCObject := serverResponse[resultDescKey.Parent]
		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
		newParentResultDescObject := make(map[string]interface{})
		switch valueOfVariableParentRC.Kind() {
		case reflect.Map:
			iter := valueOfVariableParentRC.MapRange()
			for iter.Next() {
				if iter.Key().String() == resultDescKey.Field {
					newParentResultDescObject[iter.Key().String()] = rcConfig.DescEng.String
				} else {
					newParentResultDescObject[iter.Key().String()] = iter.Value().Interface()
				}
			}
		}
		serverResponse[resultDescKey.Parent] = newParentResultDescObject
	}

	return
}
