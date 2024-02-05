package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/function"
	"alfamart-channel/helper"
	"alfamart-channel/models"
	"alfamart-channel/repo"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/zenmuharom/zenlogger"
)

type DefaultRoute struct {
	logger              zenlogger.Zenlogger
	endpoint            string
	request             models.Request
	rawRequest          map[string]any
	servRespConf        []domain.ServerResponse
	assigner            function.Assigner
	serverReqFields     map[string]models.Variable
	middlewareReqFields map[string]models.Variable
	serverResFields     map[string]models.Variable
	hasher              helper.Hasher
}

type Route interface {
	Process() (response map[string]interface{}, errRes error)
	// AssignServerResponse(productCode, endpoint string, middlewResponseIDs []int, middlwareResponse map[string]interface{}) (serverResponse map[string]interface{}, err error)
	AssignServerResponse2(productCode, endpoint string, middlewResponseIDs []int, middlwareResponse map[string]interface{}) (serverResponse map[string]interface{}, err error)
	// AssignMiddlewareRequest(middleware string, process string, userProductConf *domain.UserProduct, serverRequest map[string]interface{}) (middlewareRequest map[string]interface{}, err error)
	AssignMiddlewareRequest2(middleware string, process string, userProductConf *domain.UserProduct, serverRequest map[string]interface{}) (middlewareRequest map[string]interface{}, err error)
	OrderingServerResponse(unorderedServerResponse map[string]interface{}) (orderedServerResponse string, err error)
}

func NewRoute(logger zenlogger.Zenlogger, endpoint string, req models.Request, rawReq map[string]any) Route {

	serverResponseRepo := repo.NewServerResponseRepo(logger)
	serverResponseConfs, err := serverResponseRepo.FindAllByProductCodeAndEndpoint(req.ProductID, endpoint)
	if err != nil {
		logger.Error("NewRoute", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return &DefaultRoute{}
	}

	mrFields := make(map[string]models.Variable, 0)
	srFields := make(map[string]models.Variable, 0)

	return &DefaultRoute{
		logger:              logger,
		endpoint:            endpoint,
		request:             req,
		rawRequest:          rawReq,
		servRespConf:        serverResponseConfs,
		assigner:            function.NewAssigner(logger),
		middlewareReqFields: mrFields,
		serverResFields:     srFields,
	}
}

func (service *DefaultRoute) Process() (response map[string]interface{}, errRes error) {
	service.logger.Debug("Process", zenlogger.ZenField{Key: "endpoint", Value: service.endpoint}, zenlogger.ZenField{Key: "rawRequest", Value: service.rawRequest}, zenlogger.ZenField{Key: "request", Value: service.request})

	userProductRepo := repo.NewUserProductRepo(service.logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: service.request.AgentID}, ProductCode: sql.NullString{String: service.request.ProductID}})
	if err != nil {
		service.logger.Error("Process", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	mReq, err := service.AssignMiddlewareRequest2("ts_adapter", "inquiry", &userProductConf, service.rawRequest)
	if err != nil {
		service.logger.Error("Process", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	fmt.Println(fmt.Printf("middlewareReq: %#v", mReq))

	return
}

// func (service *DefaultRoute) AssignMiddlewareRequest(middleware string, process string, userProductConf *domain.UserProduct, serverRequest map[string]interface{}) (middlewareRequest map[string]interface{}, err error) {
// 	service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "middleware", Value: middleware}, zenlogger.ZenField{Key: "process", Value: process}, zenlogger.ZenField{Key: "serverRequest", Value: serverRequest})

// 	// set productCode
// 	productCode := userProductConf.ProductCode.String
// 	// override if product code mapped not empty
// 	if userProductConf.ProductCodeMapped.String != "" {
// 		productCode = userProductConf.ProductCodeMapped.String
// 	}

// 	// Find all middleware request
// 	middlewareRequestRepo := repo.NewMiddlewareRequestRepo(service.logger)
// 	middlewareRequestConfs, err := middlewareRequestRepo.FindAllByAttributes(domain.MiddlewareRequest{
// 		Middleware:  sql.NullString{String: middleware},
// 		Process:     sql.NullString{String: process},
// 		ProductCode: sql.NullString{String: productCode},
// 	})
// 	service.logger.Debug("middlewareRequestConfs", zenlogger.ZenField{Key: "data", Value: middlewareRequestConfs})

// 	// gather all server response fieldIDs
// 	fieldIDs := make([]int, 0)
// 	for _, middlewareRequest := range middlewareRequestConfs {
// 		fieldIDs = append(fieldIDs, middlewareRequest.Id)
// 	}

// 	// Find all middleware request value config
// 	mrvRepo := repo.NewMiddlewareRequestValueRepo(service.logger)
// 	// mrvConfs, err := mrvRepo.FindAllByAttributes(domain.MiddlewareRequestValue{CaCode: sql.NullString{String: caCode}})
// 	mrvConfs, err := mrvRepo.FindAllByFieldIDs(fieldIDs)
// 	if err != nil {
// 		service.logger.Error(err.Error())
// 		return
// 	}

// 	if err != nil {
// 		service.logger.Error(err.Error())
// 		return
// 	}

// 	assigner := function.NewAssigner(service.logger)
// 	requestValues, errs := assigner.MiddlewareRequestParse(mrvConfs, serverRequest)
// 	if len(errs) > 0 {
// 		err = errors.New("Occur some errors")
// 		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
// 		return
// 	} else {
// 		service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "requestValues", Value: requestValues})
// 	}

// 	// FOR DEBUG in local
// 	// fmt.Println(service.logger.Debug("sebelum diconsttruct", zenlogger.ZenField{Key: "requestValues", Value: requestValues}))

// 	middlewareRequest, errs = assigner.MiddlewareRequestConstruct(middlewareRequestConfs, requestValues)
// 	if len(errs) > 0 {
// 		err = errors.New("Occur some errors")
// 		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
// 		return
// 	} else {
// 		service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "middlewareRequest", Value: middlewareRequest})
// 	}

// 	// find bit18
// 	bit18Key, bit18IValue := tool.FindFieldAs(service.logger, domain.MIDDLEWARE_REQUEST, fmt.Sprintf("ts_adapter|%v", process), "bit18", middlewareRequest)
// 	if bit18Key.Parent == "" && bit18Key.Field == "" {
// 		err = errors.New("bit18 route not set yet.")
// 		service.logger.Debug(err.Error())
// 		return
// 	}

// 	// assign bit18 based on Config @ user_product
// 	if bit18IValue == nil || bit18IValue == "" {
// 		if bit18Key.Parent == "" {
// 			middlewareRequest[bit18Key.Field] = userProductConf.Bit18.String
// 		} else {
// 			// if it has parent key
// 			parentRCObject := middlewareRequest[bit18Key.Parent]
// 			valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
// 			newParentRCObject := make(map[string]interface{})
// 			switch valueOfVariableParentRC.Kind() {
// 			case reflect.Map:
// 				iter := valueOfVariableParentRC.MapRange()
// 				for iter.Next() {
// 					if iter.Key().String() == bit18Key.Field {
// 						newParentRCObject[iter.Key().String()] = userProductConf.Bit18.String
// 					} else {
// 						newParentRCObject[iter.Key().String()] = iter.Value().Interface()
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// find bit32
// 	bit32Key, bit32IValue := tool.FindFieldAs(service.logger, domain.MIDDLEWARE_REQUEST, fmt.Sprintf("ts_adapter|%v", process), "bit32", middlewareRequest)
// 	if bit32Key.Parent == "" && bit32Key.Field == "" {
// 		err = errors.New("bit32 route not set yet.")
// 		service.logger.Debug(err.Error())
// 		return
// 	}

// 	// assign bit32 based on Config @ user_product
// 	if bit32IValue == nil || bit32IValue == "" {
// 		if bit32Key.Parent == "" {
// 			middlewareRequest[bit32Key.Field] = userProductConf.Bit32.String
// 		} else {
// 			// if it has parent key
// 			parentRCObject := middlewareRequest[bit32Key.Parent]
// 			valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
// 			newParentRCObject := make(map[string]interface{})
// 			switch valueOfVariableParentRC.Kind() {
// 			case reflect.Map:
// 				iter := valueOfVariableParentRC.MapRange()
// 				for iter.Next() {
// 					if iter.Key().String() == bit32Key.Field {
// 						newParentRCObject[iter.Key().String()] = userProductConf.Bit32.String
// 					} else {
// 						newParentRCObject[iter.Key().String()] = iter.Value().Interface()
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// find bit33
// 	bit33Key, bit33IValue := tool.FindFieldAs(service.logger, domain.MIDDLEWARE_REQUEST, fmt.Sprintf("ts_adapter|%v", process), "bit33", middlewareRequest)
// 	if bit33Key.Parent == "" && bit33Key.Field == "" {
// 		err = errors.New("bit33 route not set yet.")
// 		service.logger.Debug(err.Error())
// 		return
// 	}

// 	// assign bit33 based on Config @ user_product
// 	if bit33IValue == nil || bit33IValue == "" {
// 		if bit33Key.Parent == "" {
// 			middlewareRequest[bit33Key.Field] = userProductConf.Bit33.String
// 		} else {
// 			// if it has parent key
// 			parentRCObject := middlewareRequest[bit33Key.Parent]
// 			valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
// 			newParentRCObject := make(map[string]interface{})
// 			switch valueOfVariableParentRC.Kind() {
// 			case reflect.Map:
// 				iter := valueOfVariableParentRC.MapRange()
// 				for iter.Next() {
// 					if iter.Key().String() == bit33Key.Field {
// 						newParentRCObject[iter.Key().String()] = userProductConf.Bit33.String
// 					} else {
// 						newParentRCObject[iter.Key().String()] = iter.Value().Interface()
// 					}
// 				}
// 			}
// 		}
// 	}

// 	service.logger.Debug("AssignMiddlewareRequest", zenlogger.ZenField{Key: "middlewareRequest", Value: middlewareRequest}, zenlogger.ZenField{Key: "requestValues", Value: requestValues})

// 	return
// }

func (service *DefaultRoute) AssignServerRequest2(endpoint string, process string, userProductConf *domain.UserProduct, serverRequest map[string]interface{}) (middlewareRequest map[string]interface{}, err error) {
	service.logger.Debug("AssignServerRequest2",
		zenlogger.ZenField{Key: "serverReqFields", Value: service.serverReqFields},
	)

	// Find all middleware request
	serverRequestRepo := repo.NewServerRequestRepo(service.logger)
	srFields, err := serverRequestRepo.FindAllByEndpoint(endpoint)

	// store middleware request to memory
	err = service.StoreServerReqToMemory(srFields)
	if err != nil {
		service.logger.Error("AssignServerRequest2", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addon", Value: "error while call StoreFieldToMemory"})
	}

	// Select field value config
	fieldIDs := make([]int, 0)
	for _, field := range srFields {
		fieldIDs = append(fieldIDs, field.Id)
	}

	service.logger.Debug("AssignServerRequest2",
		zenlogger.ZenField{Key: "middlewareReqFields", Value: service.middlewareReqFields},
	)

	return
}

func (service *DefaultRoute) AssignMiddlewareRequest2(middleware string, process string, userProductConf *domain.UserProduct, serverRequest map[string]interface{}) (middlewareRequest map[string]interface{}, err error) {
	service.logger.Debug("AssignMiddlewareRequest2", zenlogger.ZenField{Key: "middleware", Value: middleware}, zenlogger.ZenField{Key: "process", Value: process}, zenlogger.ZenField{Key: "serverRequest", Value: serverRequest})

	// set productCode
	productCode := userProductConf.ProductCode.String
	// override if product code mapped not empty
	if userProductConf.ProductCodeMapped.String != "" {
		productCode = userProductConf.ProductCodeMapped.String
	}

	// Find all middleware request
	middlewareRequestRepo := repo.NewMiddlewareRequestRepo(service.logger)
	mrFields, err := middlewareRequestRepo.FindAllByAttributes(domain.MiddlewareRequest{
		Middleware:  sql.NullString{String: middleware},
		Process:     sql.NullString{String: process},
		ProductCode: sql.NullString{String: productCode},
	})

	// store middleware request to memory
	err = service.StoreMiddlewareReqToMemory(mrFields)
	if err != nil {
		service.logger.Error("AssignMiddlewareRequest2", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addon", Value: "error while call StoreFieldToMemory"})
	}

	// Select field value config
	fieldIDs := make([]int, 0)
	for _, field := range mrFields {
		fieldIDs = append(fieldIDs, field.Id)
	}
	fieldValueRepo := repo.NewMiddlewareRequestValueRepo(service.logger)
	fieldValues, err := fieldValueRepo.FindAllByFieldIDs(fieldIDs)
	if err != nil {
		service.logger.Error("AssignMiddlewareRequest2", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call productServiceRepo.FindByServiceAliasAndProductCode() function"})
		return
	}

	// assign value
	_, err = service.AssignMiddleReqValuesBasedOnConfig(fieldValues)
	if err != nil {
		service.logger.Error("AssignMiddlewareRequest2", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call AssignValuesBasedOnConfig"})
		return
	}

	service.logger.Debug("AssignMiddlewareRequest2",
		zenlogger.ZenField{Key: "middlewareReqFields", Value: service.middlewareReqFields},
	)

	return
}

func (service *DefaultRoute) StoreServerReqToMemory(fields []domain.ServerRequest) (errRes error) {

	service.logger.Debug("StoreServerReqToMemory",
		zenlogger.ZenField{Key: "fields", Value: fields},
		zenlogger.ZenField{Key: "middlewareReqFields", Value: service.middlewareReqFields},
	)

	for index, field := range fields {
		service.logger.Debug("StoreServerReqToMemory", zenlogger.ZenField{Key: "index", Value: index}, zenlogger.ZenField{Key: "field", Value: field})

		variable := models.Variable{
			Field: field.Field.String,
			Type:  field.Type.String,
			As:    field.FieldAs.String,
			Value: nil,
		}

		if field.ParentId.Valid {
			variable.Parent = service.hasher.GenerateHash(service.logger.GetPid(), field.ParentId.Int64)
		}

		memoryKey := service.hasher.GenerateHash(service.logger.GetPid(), field.Id)
		service.middlewareReqFields[memoryKey] = variable

		// assign to it's parent
		if field.ParentId.Valid {
			parentMK := service.hasher.GenerateHash(service.logger.GetPid(), field.ParentId.Int64)
			parentNode := service.middlewareReqFields[parentMK]
			// if parent is already exist
			if _, ok := service.middlewareReqFields[parentMK]; ok {

				parentNode = service.middlewareReqFields[parentMK]

				// if parent value is nil
				if parentNode.Childs == nil {
					parentNode.Childs = map[string]models.Variable{
						memoryKey: variable,
					}
				} else {
					// if parent already have child
					parentNode.Childs[memoryKey] = variable
				}

				// reassign back
				service.middlewareReqFields[parentMK] = parentNode

			} else {
				// if parent not exist yet
				service.middlewareReqFields[parentMK] = models.Variable{
					Field:   service.middlewareReqFields[parentMK].Field,
					Type:    service.middlewareReqFields[parentMK].Type,
					As:      service.middlewareReqFields[parentMK].As,
					Parent:  service.middlewareReqFields[parentMK].Parent,
					RouteId: service.middlewareReqFields[parentMK].RouteId,
					Childs: map[string]models.Variable{
						memoryKey: variable,
					},
				}
			}
		}
	}

	service.logger.Debug("StoreServerReqToMemory",
		zenlogger.ZenField{Key: "fields", Value: fields},
		zenlogger.ZenField{Key: "serverReqFields", Value: service.serverReqFields},
	)

	return
}

func (service *DefaultRoute) StoreMiddlewareReqToMemory(fields []domain.MiddlewareRequest) (errRes error) {

	service.logger.Debug("StoreMiddlewareReqToMemory",
		zenlogger.ZenField{Key: "fields", Value: fields},
		zenlogger.ZenField{Key: "middlewareReqFields", Value: service.middlewareReqFields},
	)

	for index, field := range fields {
		service.logger.Debug("StoreMiddlewareReqToMemory", zenlogger.ZenField{Key: "index", Value: index}, zenlogger.ZenField{Key: "field", Value: field})

		variable := models.Variable{
			Field: field.Field.String,
			Type:  field.Type.String,
			As:    field.FieldAs.String,
			Value: nil,
		}

		if field.ParentId.Valid {
			variable.Parent = service.hasher.GenerateHash(service.logger.GetPid(), field.ParentId.Int64)
		}

		memoryKey := service.hasher.GenerateHash(service.logger.GetPid(), field.Id)
		service.middlewareReqFields[memoryKey] = variable

		// assign to it's parent
		if field.ParentId.Valid {
			parentMK := service.hasher.GenerateHash(service.logger.GetPid(), field.ParentId.Int64)
			parentNode := service.middlewareReqFields[parentMK]
			// if parent is already exist
			if _, ok := service.middlewareReqFields[parentMK]; ok {

				parentNode = service.middlewareReqFields[parentMK]

				// if parent value is nil
				if parentNode.Childs == nil {
					parentNode.Childs = map[string]models.Variable{
						memoryKey: variable,
					}
				} else {
					// if parent already have child
					parentNode.Childs[memoryKey] = variable
				}

				// reassign back
				service.middlewareReqFields[parentMK] = parentNode

			} else {
				// if parent not exist yet
				service.middlewareReqFields[parentMK] = models.Variable{
					Field:   service.middlewareReqFields[parentMK].Field,
					Type:    service.middlewareReqFields[parentMK].Type,
					As:      service.middlewareReqFields[parentMK].As,
					Parent:  service.middlewareReqFields[parentMK].Parent,
					RouteId: service.middlewareReqFields[parentMK].RouteId,
					Childs: map[string]models.Variable{
						memoryKey: variable,
					},
				}
			}
		}
	}

	service.logger.Debug("StoreMiddlewareReqToMemory",
		zenlogger.ZenField{Key: "fields", Value: fields},
		zenlogger.ZenField{Key: "middlewareReqFields", Value: service.middlewareReqFields},
	)

	return
}

func (service *DefaultRoute) AssignMiddleReqValuesBasedOnConfig(middlewareRequestValues []domain.MiddlewareRequestValue) (filledFields map[string]models.Variable, errRes error) {

	service.logger.Debug("AssignMiddleReqValuesBasedOnConfig",
		zenlogger.ZenField{Key: "referenceFields", Value: service.middlewareReqFields},
		zenlogger.ZenField{Key: "middlewareRequestValues", Value: middlewareRequestValues},
	)

	// TODO
	fieldValues := service.convertMiddlewareRequestValueToFieldValue(middlewareRequestValues)

	for index, fieldValue := range fieldValues {
		service.logger.Debug("AssignMiddleReqValuesBasedOnConfig", zenlogger.ZenField{Key: "index", Value: index}, zenlogger.ZenField{Key: "fieldValue", Value: fieldValue})
		value := service.assigner.AssignMiddlewareRequestValue(fieldValue, service.serverReqFields)

		memoryKey := service.hasher.GenerateHash(service.logger.GetPid(), fieldValue.FieldId)
		currentNode := service.middlewareReqFields[memoryKey]
		currentNode.Value = value
		service.middlewareReqFields[memoryKey] = currentNode
		if currentNode.Parent != "" {
			service.UpdateChildValueInParent(currentNode.Parent, memoryKey, value, service.middlewareReqFields)
		}

	}
	service.logger.Debug("AssignMiddleReqValuesBasedOnConfig",
		zenlogger.ZenField{Key: "referenceFields", Value: service.middlewareReqFields},
		zenlogger.ZenField{Key: "middlewareRequestValues", Value: middlewareRequestValues},
	)

	return
}

func (service *DefaultRoute) convertMiddlewareRequestToVariable(mrFields []domain.MiddlewareRequest) (fields []models.Variable) {
	service.logger.Debug("convertMiddlewareRequestToVariable", zenlogger.ZenField{Key: "mrFields", Value: mrFields})
	// TODO

	return
}

func (service *DefaultRoute) convertMiddlewareRequestValueToFieldValue(mrValues []domain.MiddlewareRequestValue) (fieldValues []domain.FieldValue) {

	service.logger.Debug("convertMiddlewareRequestValueToFieldValue", zenlogger.ZenField{Key: "mrValues", Value: mrValues}, zenlogger.ZenField{Key: "fieldValues", Value: fieldValues})

	for _, mrValue := range mrValues {

		fieldValue := domain.FieldValue{
			Id:                mrValue.Id,
			FieldId:           mrValue.FieldId.Int64,
			FieldName:         mrValue.FieldName,
			FieldType:         mrValue.FieldType.String,
			FieldAs:           mrValue.FieldAs,
			FieldParentId:     mrValue.ParentId,
			FieldRefId:        mrValue.ServerRequestId,
			Args:              mrValue.Args,
			ConditionFieldId:  mrValue.ConditionFieldId,
			ConditionOperator: mrValue.ConditionOperator,
			ConditionValue:    mrValue.ConditionOperator,
		}

		service.logger.Debug("convertMiddlewareRequestValueToFieldValue", zenlogger.ZenField{Key: "mrValue", Value: mrValue}, zenlogger.ZenField{Key: "fieldValues", Value: fieldValues})

		fieldValues = append(fieldValues, fieldValue)
	}

	service.logger.Debug("convertMiddlewareRequestValueToFieldValue", zenlogger.ZenField{Key: "mrValues", Value: mrValues}, zenlogger.ZenField{Key: "fieldValues", Value: fieldValues})

	return
}

// Update the Value of a child Variable in a parent Variable and also update the grandparent if needed
func (service *DefaultRoute) UpdateChildValueInParent(parentKey, childKey string, newValue interface{}, fields map[string]models.Variable) {
	if parentNode, parentExists := fields[parentKey]; parentExists {
		if childNode, childExists := parentNode.Childs[childKey]; childExists {
			childNode.Value = newValue
			parentNode.Childs[childKey] = childNode
			fields[parentKey] = parentNode
			// If the parent has a parent (grandparent), update the grandparent as well
			if parentNode.Parent != "" {
				service.UpdateChildValueInParent(parentNode.Parent, parentKey, newValue, fields)
			}
		}
	}
}

// func (service *DefaultRoute) AssignServerResponse(productCode, endpoint string, middlewResponseIDs []int, middlewareResponse map[string]interface{}) (serverResponse map[string]interface{}, err error) {
// 	service.logger.Debug("AssignServerResponse", zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "middlewareResponse", Value: middlewareResponse})

// 	var serverResponseConfs []domain.ServerResponse

// 	// find all server response
// 	serverResponseRepo := repo.NewServerResponseRepo(service.logger)
// 	// if productCode is not empty select it base on productCode otherwise select the default one
// 	serverResponseConfs, err = serverResponseRepo.FindAllByProductCodeAndEndpoint(productCode, endpoint)
// 	service.logger.Debug("serverResponseConfs", zenlogger.ZenField{Key: "data", Value: serverResponseConfs})
// 	if err != nil {
// 		service.logger.Error(err.Error())
// 		return
// 	}

// 	if len(serverResponseConfs) == 0 {
// 		err = errors.New("There is no server response config found")
// 		service.logger.Info("AssignServerResponse", zenlogger.ZenField{Key: "error", Value: err.Error()})
// 		return
// 	}

// 	// gather all server response fieldIDs
// 	fieldIDs := make([]int, 0)
// 	for _, serverResponse := range serverResponseConfs {
// 		fieldIDs = append(fieldIDs, serverResponse.Id)
// 	}

// 	// Find all server response value config
// 	srvRepo := repo.NewServerResponseValueRepo(service.logger)
// 	// srvConfs, err := srvRepo.FindAllByProductCode(productCode)
// 	srvConfs, err := srvRepo.FindAllByServerResponseID(productCode, fieldIDs, middlewResponseIDs)
// 	if err != nil {
// 		service.logger.Error(err.Error())
// 		return
// 	}

// 	assigner := function.NewAssigner(service.logger)

// 	responseValues, errs := assigner.ServerResponseParse(srvConfs, middlewareResponse)
// 	if len(errs) > 0 {
// 		err = errors.New("Occur some errors")
// 		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
// 		return
// 	}

// 	serverResponse, errs = assigner.ServerResponseConstruct(serverResponseConfs, responseValues)
// 	if len(errs) > 0 {
// 		err = errors.New("Occur some errors")
// 		service.logger.Error("Occur some errors", zenlogger.ZenField{Key: "errors", Value: fmt.Sprintf("%#v", errs)})
// 		return
// 	}

// 	// find rc key
// 	resultCodeKey, resultCodeIValue := tool.FindFieldAs(service.logger, domain.SERVER_RESPONSE, endpoint, "resultCode", serverResponse)
// 	rcConfigRepo := repo.NewRcConfigRepo(service.logger)
// 	rcCode, _ := strconv.ParseInt(fmt.Sprintf("%v", resultCodeIValue), 10, 64)
// 	rcConfig, err := rcConfigRepo.FindRc(productCode, rcCode)

// 	if resultCodeKey.Parent == "" {
// 		serverResponse[resultCodeKey.Field] = rcConfig.Code.String
// 	} else {
// 		// if it has parent key
// 		parentRCObject := serverResponse[resultCodeKey.Parent]
// 		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
// 		newParentResultCodeObject := make(map[string]interface{})
// 		switch valueOfVariableParentRC.Kind() {
// 		case reflect.Map:
// 			iter := valueOfVariableParentRC.MapRange()
// 			for iter.Next() {
// 				if iter.Key().String() == resultCodeKey.Field {
// 					newParentResultCodeObject[iter.Key().String()] = rcConfig.Code.String
// 				} else {
// 					newParentResultCodeObject[iter.Key().String()] = iter.Value().Interface()
// 				}
// 			}
// 		}
// 		serverResponse[resultCodeKey.Parent] = newParentResultCodeObject
// 	}

// 	resultDescKey, _ := tool.FindFieldAs(service.logger, domain.SERVER_RESPONSE, endpoint, "resultDesc", serverResponse)
// 	if resultDescKey.Parent == "" {
// 		serverResponse[resultDescKey.Field] = rcConfig.DescId.String
// 	} else {
// 		// if it has parent key
// 		parentRCObject := serverResponse[resultDescKey.Parent]
// 		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
// 		newParentResultDescObject := make(map[string]interface{})
// 		switch valueOfVariableParentRC.Kind() {
// 		case reflect.Map:
// 			iter := valueOfVariableParentRC.MapRange()
// 			for iter.Next() {
// 				if iter.Key().String() == resultDescKey.Field {
// 					newParentResultDescObject[iter.Key().String()] = rcConfig.DescId.String
// 				} else {
// 					newParentResultDescObject[iter.Key().String()] = iter.Value().Interface()
// 				}
// 			}
// 		}
// 		serverResponse[resultDescKey.Parent] = newParentResultDescObject
// 	}

// 	return
// }

func (service *DefaultRoute) StoreServerResToMemory(fields []domain.ServerResponse) (errRes error) {

	service.logger.Debug("StoreServerResToMemory",
		zenlogger.ZenField{Key: "fields", Value: fields},
	)

	for index, field := range fields {
		service.logger.Debug("StoreServerResToMemory", zenlogger.ZenField{Key: "index", Value: index}, zenlogger.ZenField{Key: "field", Value: field})

		variable := models.Variable{
			Field: field.Field.String,
			Type:  field.Type.String,
			As:    field.FieldAs.String,
			Value: nil,
		}

		if field.ParentId.Valid {
			variable.Parent = service.hasher.GenerateHash(service.logger.GetPid(), field.ParentId.Int64)
		}

		memoryKey := service.hasher.GenerateHash(service.logger.GetPid(), field.Id)
		service.serverResFields[memoryKey] = variable

		// assign to it's parent
		if field.ParentId.Valid {
			parentMK := service.hasher.GenerateHash(service.logger.GetPid(), field.ParentId.Int64)
			parentNode := service.serverResFields[parentMK]
			// if parent is already exist
			if _, ok := service.serverResFields[parentMK]; ok {

				parentNode = service.serverResFields[parentMK]

				// if parent value is nil
				if parentNode.Childs == nil {
					parentNode.Childs = map[string]models.Variable{
						memoryKey: variable,
					}
				} else {
					// if parent already have child
					parentNode.Childs[memoryKey] = variable
				}

				// reassign back
				service.serverResFields[parentMK] = parentNode

			} else {
				// if parent not exist yet
				service.serverResFields[parentMK] = models.Variable{
					Field:  service.serverResFields[parentMK].Field,
					Type:   service.serverResFields[parentMK].Type,
					As:     service.serverResFields[parentMK].As,
					Parent: service.serverResFields[parentMK].Parent,
					Childs: map[string]models.Variable{
						memoryKey: variable,
					},
				}
			}
		}
	}

	service.logger.Debug("StoreServerResToMemory",
		zenlogger.ZenField{Key: "fields", Value: fields},
	)

	return
}

func (service *DefaultRoute) AssignServerResponse2(productCode, endpoint string, middlewResponseIDs []int, middlewareResponse map[string]interface{}) (serverResponse map[string]interface{}, err error) {
	service.logger.Debug("AssignServerResponse", zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "middlewareResponse", Value: middlewareResponse})

	var serverResponseConfs []domain.ServerResponse

	// find all server response
	serverResponseRepo := repo.NewServerResponseRepo(service.logger)
	// if productCode is not empty select it base on productCode otherwise select the default one
	serverResponseConfs, err = serverResponseRepo.FindAllByProductCodeAndEndpoint(productCode, endpoint)
	service.logger.Debug("AssignServerResponse2", zenlogger.ZenField{Key: "data", Value: serverResponseConfs})
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	if len(serverResponseConfs) == 0 {
		err = errors.New("There is no server response config found")
		service.logger.Info("AssignServerResponse", zenlogger.ZenField{Key: "error", Value: err.Error()})
		return
	}

	// store middleware request to memory
	err = service.StoreServerResToMemory(serverResponseConfs)
	if err != nil {
		service.logger.Error("AssignServerResponse2", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addon", Value: "error while call StoreFieldToMemory"})
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

	_, err = service.AssignServerResValuesBasedOnConfig(srvConfs)
	if err != nil {
		service.logger.Error("AssignServerResponse2", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call AssignValuesBasedOnConfig"})
		return
	}

	return
}

func (service *DefaultRoute) convertServerResponseValueToFieldValue(srValues []domain.ServerResponseValue) (fieldValues []domain.FieldValue) {

	service.logger.Debug("convertServerResponseValueToFieldValue", zenlogger.ZenField{Key: "srValues", Value: srValues}, zenlogger.ZenField{Key: "fieldValues", Value: fieldValues})

	for _, srValue := range srValues {

		fieldValue := domain.FieldValue{
			Id:                srValue.Id,
			FieldId:           srValue.FieldId.Int64,
			FieldName:         srValue.FieldName,
			FieldType:         srValue.FieldType.String,
			FieldAs:           srValue.FieldAs,
			FieldParentId:     srValue.ParentId,
			FieldRefId:        srValue.MiddlewareResponseId,
			Args:              srValue.Args,
			ConditionFieldId:  srValue.ConditionFieldId,
			ConditionOperator: srValue.ConditionOperator,
			ConditionValue:    srValue.ConditionOperator,
		}

		service.logger.Debug("convertServerResponseValueToFieldValue", zenlogger.ZenField{Key: "srValue", Value: srValue}, zenlogger.ZenField{Key: "fieldValues", Value: fieldValues})

		fieldValues = append(fieldValues, fieldValue)
	}

	service.logger.Debug("convertServerResponseValueToFieldValue", zenlogger.ZenField{Key: "srValues", Value: srValues}, zenlogger.ZenField{Key: "fieldValues", Value: fieldValues})

	return

}

func (service *DefaultRoute) AssignServerResValuesBasedOnConfig(srValues []domain.ServerResponseValue) (filledFields map[string]models.Variable, errRes error) {

	service.logger.Debug("AssignServerResValuesBasedOnConfig",
		zenlogger.ZenField{Key: "referenceFields", Value: service.serverResFields},
		zenlogger.ZenField{Key: "srValues", Value: srValues},
	)

	// TODO
	// convert server response value to field value
	fieldValues := service.convertServerResponseValueToFieldValue(srValues)

	for index, fieldValue := range fieldValues {
		service.logger.Debug("AssignServerResValuesBasedOnConfig", zenlogger.ZenField{Key: "index", Value: index}, zenlogger.ZenField{Key: "fieldValue", Value: fieldValue})
		value := service.assigner.AssignServerResponseValue(fieldValue, service.serverResFields)

		memoryKey := service.hasher.GenerateHash(service.logger.GetPid(), fieldValue.FieldId)
		currentNode := service.serverResFields[memoryKey]
		currentNode.Value = value
		service.serverResFields[memoryKey] = currentNode
		if currentNode.Parent != "" {
			service.UpdateChildValueInParent(currentNode.Parent, memoryKey, value, service.serverResFields)
		}

	}

	service.logger.Debug("AssignServerResValuesBasedOnConfig",
		zenlogger.ZenField{Key: "referenceFields", Value: service.serverResFields},
		zenlogger.ZenField{Key: "fieldValues", Value: fieldValues},
	)

	return
}

func (service *DefaultRoute) OrderingServerResponse(unorderedServerResponse map[string]interface{}) (orderedServerResponse string, err error) {

	service.logger.Debug("OrderingServerResponse", zenlogger.ZenField{Key: "unorderedServerResponse", Value: unorderedServerResponse})

	var serverResponseConstruct = make([]string, 0)

	sort.Slice(service.servRespConf, func(i, j int) bool {
		return service.servRespConf[i].Order < service.servRespConf[j].Order
	})

	for _, srvResp := range service.servRespConf {
		serverResponseConstruct = append(serverResponseConstruct, fmt.Sprintf("%v", unorderedServerResponse[srvResp.Field.String]))
	}

	orderedServerResponse = strings.Join(serverResponseConstruct, "|")

	return
}
