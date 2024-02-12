package api

import (
	"alfamart-channel/domain"
	"alfamart-channel/models"
	"alfamart-channel/repo"
	"alfamart-channel/service"
	"alfamart-channel/tool"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

func (server *DefaultServer) validate() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		endpoint := ctx.Request.URL.Path

		logger := server.setupLogger()
		// read parameter URI
		// Get the map of query parameters from the URL
		req := make(map[string]interface{})
		queryParams := ctx.Request.URL.Query()
		// Iterate through the map and print the key-value pairs
		for key, values := range queryParams {
			// If there is only one value, add it directly to the map
			if len(values) == 1 {
				req[key] = values[0]
			} else {
				// If there are multiple values, add them as a slice of strings
				req[key] = values
			}
		}
		logger.Debug("validate", zenlogger.ZenField{Key: "queryParams", Value: queryParams}, zenlogger.ZenField{Key: "req", Value: req})

		// find productCode
		_, productCodeIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "productCode", req)
		if productCodeIValue == nil {
			err := errors.New("Product code route not set yet")
			logger.Error(err.Error())
			errorMsg := models.ErrorMsg{GoCode: 1706, Err: err}
			logger.Debug("sendErrorResponse", zenlogger.ZenField{Key: "requestURI", Value: ctx.Request.URL.Path})
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, "", errorMsg)
			return
		}
		productCode := fmt.Sprintf("%v", productCodeIValue)
		logger.Debug("validate", zenlogger.ZenField{Key: "productCode", Value: productCode})

		// validate request body using configuration DB
		if valid, errRes := validateRequest(logger, ctx.Request.URL.Path, req); !valid {
			logger.Error("validate", zenlogger.ZenField{Key: "valid", Value: valid})
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, ctx.Request.URL.Path, errRes)
			return
		}

		ipConfigRepo := repo.NewIpConfigRepo(logger)

		valid, err := ipConfigRepo.FindIp(ctx.ClientIP())
		if err != nil {
			logger.Error(err.Error())
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, models.ErrorMsg{GoCode: 1706, Err: err})
			return
		}

		if !valid {
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, models.ErrorMsg{GoCode: 1994, Err: err})
			return
		}

		signatureService := service.NewSignatureService(logger)
		valid, err = signatureService.Check(ctx.Request.URL.Path, req)
		if err != nil {
			logger.Error(err.Error())
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, models.ErrorMsg{GoCode: 1706, Err: err})
			return
		}
		if !valid {
			logger.Info("Invalid signature")
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, models.ErrorMsg{GoCode: 36, Err: err})
			return
		}

		// newBody, err := json.Marshal(req)
		// if err != nil {
		// 	logger.Error("Fail to re-wrap request", zenlogger.ZenField{Key: "error", Value: err})
		// 	sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, models.ErrorMsg{GoCode: 7000, Err: err})
		// 	return
		// }
		// ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))

		// ctx.Next()
	}
}

func validateRequest(logger zenlogger.Zenlogger, endpoint string, req map[string]interface{}) (valid bool, errRes models.ErrorMsg) {
	valid = true
	logger.Debug("validateRequest", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "req", Value: req})
	validationMsg := ""

	serverRequestRepo := repo.NewServerRequestRepo(logger)
	configs, err := serverRequestRepo.FindRequestQueryByEndpoint(endpoint)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Debug("validateRequest", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "req", Value: req}, zenlogger.ZenField{Key: "configs", Value: configs})
	for _, field := range configs {
		value := ""
		ok := false

		_, ok = req[field.Field.String]
		// if field exist check if it is not empty
		if ok {
			ok = !(req[field.Field.String] == "" || req[field.Field.String] == "null" || req[field.Field.String] == nil)
			value = fmt.Sprintf("%v", req[field.Field.String])
		}

		if field.Required.Bool && !ok {
			if valid {
				errRes.GoCode = 35
				valid = false
			} else {
				validationMsg += ", "
			}
			validationMsg += fmt.Sprintf("%s is missing", field.Field.String)
		}

		if ok {
			// check the minimum length validation of value
			if field.LengthMin.Valid {
				ok = (len(value) >= int(field.LengthMin.Int64))
				if !ok {
					valid = false
					validationMsg += fmt.Sprintf("%s minimum length is %v", field.Field.String, field.LengthMin.Int64)
				}
			}

			// check the maximum length validation of value
			if field.LengthMax.Valid {
				ok = (len(value) <= int(field.LengthMax.Int64))
				if !ok {
					valid = false
					validationMsg += fmt.Sprintf("%s maximum length is %v", field.Field.String, field.LengthMax.Int64)
				}
			}
		}

		logger.Debug("validateRequest", zenlogger.ZenField{Key: "field", Value: field.Field.String}, zenlogger.ZenField{Key: "value", Value: value}, zenlogger.ZenField{Key: "required", Value: field.Required.Bool}, zenlogger.ZenField{Key: "forbidded", Value: field.Forbidded.Bool}, zenlogger.ZenField{Key: "validationMsg", Value: validationMsg})
	}

	errRes.Err = errors.New(validationMsg)

	if !valid {
		return
	}

	return
}
