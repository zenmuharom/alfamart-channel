package tool

import (
	"alfamart-channel/domain"
	"alfamart-channel/repo"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zenmuharom/zenlogger"
)

type VariableField struct {
	Field  string
	Parent string
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "number":
		return "must a number"
	case "len":
		return "must have minimum length " + fe.Param()
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}
	return fe.Error()
}

func findFieldAsServerRequest(logger zenlogger.Zenlogger, endpoint, field string, req map[string]interface{}) (index VariableField, value interface{}) {
	serverRequestRepo := repo.NewServerRequestRepo(logger)
	config, err := serverRequestRepo.FindByEndpointAndFieldAs(endpoint, field)
	if err != nil {
		logger.Error(err.Error())
		value = nil
		return
	}

	if config.FieldParent.String != "" {
		// Assert the interface{} to a map[string]interface{}
		parentObject, ok := req[config.FieldParent.String].(map[string]interface{})
		index.Parent = config.FieldParent.String
		index.Field = config.Field.String
		if ok {
			_, ok := parentObject[config.Field.String]
			if ok {
				value = parentObject[config.Field.String]
			} else {
				value = nil
			}
		} else {
			value = nil
		}
	} else {
		index.Field = config.Field.String
		_, ok := req[config.Field.String]
		if ok {
			value = req[config.Field.String]
		}
	}

	return
}

func findFieldAsServerResponse(logger zenlogger.Zenlogger, endpoint, field string, req map[string]interface{}) (index VariableField, value interface{}) {
	serverResponseRepo := repo.NewServerResponseRepo(logger)
	config, err := serverResponseRepo.FindByEndpointAndFieldAs(endpoint, field)
	if err != nil {
		logger.Error(err.Error())
		value = nil
		return
	}

	if config.FieldParent.String != "" {
		// Assert the interface{} to a map[string]interface{}
		parentObject, ok := req[config.FieldParent.String].(map[string]interface{})
		index.Parent = config.FieldParent.String
		index.Field = config.Field.String
		if ok {
			_, ok := parentObject[config.Field.String]
			if ok {
				value = parentObject[config.Field.String]
			} else {
				value = nil
			}
		} else {
			value = nil
		}
	} else {
		index.Field = config.Field.String
		_, ok := req[config.Field.String]
		if ok {
			value = req[config.Field.String]
		}
	}

	return
}

func findFieldAsMiddlewareResponse(logger zenlogger.Zenlogger, middleware, process, field string, req map[string]interface{}) (index VariableField, value interface{}) {
	middlewareResponseRepo := repo.NewMiddlewareResponseRepo(logger)
	config, err := middlewareResponseRepo.FindByProcessAndFieldAs(middleware, process, field)
	if err != nil {
		logger.Error(err.Error())
		value = nil
		return
	}

	if config.FieldParent.String != "" {
		// Assert the interface{} to a map[string]interface{}
		parentObject, ok := req[config.FieldParent.String].(map[string]interface{})
		index.Parent = config.FieldParent.String
		index.Field = config.Field.String
		if ok {
			_, ok := parentObject[config.Field.String]
			if ok {
				value = parentObject[config.Field.String]
			} else {
				value = nil
			}
		} else {
			value = nil
		}
	} else {
		index.Field = config.Field.String
		_, ok := req[config.Field.String]
		if ok {
			value = req[config.Field.String]
		}
	}

	return
}

func findFieldAsMiddlewareRequest(logger zenlogger.Zenlogger, middleware, process, field string, req map[string]interface{}) (index VariableField, value interface{}) {
	middlewareRequestRepo := repo.NewMiddlewareRequestRepo(logger)
	config, err := middlewareRequestRepo.FindByProcessAndFieldAs(middleware, process, field)
	if err != nil {
		logger.Error(err.Error())
		value = nil
		return
	}

	if config.FieldParent.String != "" {
		// Assert the interface{} to a map[string]interface{}
		parentObject, ok := req[config.FieldParent.String].(map[string]interface{})
		index.Parent = config.FieldParent.String
		index.Field = config.Field.String
		if ok {
			_, ok := parentObject[config.Field.String]
			if ok {
				value = parentObject[config.Field.String]
			} else {
				value = nil
			}
		} else {
			value = nil
		}
	} else {
		index.Field = config.Field.String
		_, ok := req[config.Field.String]
		if ok {
			value = req[config.Field.String]
		}
	}

	return
}

func FindFieldAs(logger zenlogger.Zenlogger, table, endpoint, field string, variableObject map[string]interface{}) (index VariableField, value interface{}) {

	switch table {
	case domain.SERVER_REQUEST:
		index, value = findFieldAsServerRequest(logger, endpoint, field, variableObject)
	case domain.SERVER_RESPONSE:
		index, value = findFieldAsServerResponse(logger, endpoint, field, variableObject)
	case domain.MIDDLEWARE_REQUEST:
		middleware := ""
		process := ""

		midprocess := strings.Split(endpoint, "|")

		middleware = fmt.Sprintf("%v", midprocess[0])
		if len(midprocess) > 1 {
			process = fmt.Sprintf("%v", midprocess[1])
		}

		index, value = findFieldAsMiddlewareRequest(logger, middleware, process, field, variableObject)
	case domain.MIDDLEWARE_RESPONSE:
		middleware := ""
		process := ""

		midprocess := strings.Split(endpoint, "|")

		middleware = fmt.Sprintf("%v", midprocess[0])
		if len(midprocess) > 1 {
			process = fmt.Sprintf("%v", midprocess[1])
		}

		index, value = findFieldAsMiddlewareResponse(logger, middleware, process, field, variableObject)
	}

	return
}
