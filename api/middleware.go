package api

import (
	"alfamart-channel/domain"
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
		_, productCodeIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, ctx.Request.URL.Path, "productCode", req)
		if productCodeIValue == nil {
			err := errors.New("Product code route not set yet")
			logger.Error(err.Error())
			errorMsg := ErrorMsg{GoCode: 1706, Err: err}
			logger.Debug("sendErrorResponse", zenlogger.ZenField{Key: "requestURI", Value: ctx.Request.URL.Path})
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, "", errorMsg)
			return
		}
		productCode := fmt.Sprintf("%v", productCodeIValue)
		logger.Debug("validate", zenlogger.ZenField{Key: "productCode", Value: productCode})

		// validate request body using configuration DB
		if valid, errRes := validateRequest(logger, ctx.Request.URL.Path, req); !valid {
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, ctx.Request.URL.Path, errRes)
			return
		}

		ipConfigRepo := repo.NewIpConfigRepo(logger)

		valid, err := ipConfigRepo.FindIp(ctx.ClientIP())
		if err != nil {
			logger.Error(err.Error())
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, ErrorMsg{GoCode: 1706, Err: err})
			return
		}

		if !valid {
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, ErrorMsg{GoCode: 1994, Err: err})
			return
		}

		signatureService := service.NewSignatureService(logger)
		valid, err = signatureService.Check(ctx.Request.URL.Path, req)
		if err != nil {
			logger.Error(err.Error())
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, ErrorMsg{GoCode: 1706, Err: err})
			return
		}
		if !valid {
			logger.Info("Invalid signature")
			sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, ErrorMsg{GoCode: 36, Err: err})
			return
		}

		// newBody, err := json.Marshal(req)
		// if err != nil {
		// 	logger.Error("Fail to re-wrap request", zenlogger.ZenField{Key: "error", Value: err})
		// 	sendErrorResponse(ctx, logger, ctx.Request.URL.Path, productCode, ErrorMsg{GoCode: 7000, Err: err})
		// 	return
		// }
		// ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))

		// ctx.Next()
	}
}

func validateRequest(logger zenlogger.Zenlogger, endpoint string, req map[string]interface{}) (valid bool, errRes ErrorMsg) {
	valid = true
	logger.Debug("validateRequest", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "req", Value: req})
	validationMsg := ""

	serverRequestRepo := repo.NewServerRequestRepo(logger)
	configs, err := serverRequestRepo.FindAllByEndpoint(endpoint)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("validateRequest", zenlogger.ZenField{Key: "configs", Value: configs})
	for _, field := range configs {
		ok := false
		// check if config has parent
		if field.ParentId.Int64 == 0 {
			_, ok = req[field.Field.String]
			// if field exist check if it is not empty
			if ok {
				ok = !(req[field.Field.String] == "" || req[field.Field.String] == "null" || req[field.Field.String] == nil)
			}
		} else {
			// otherwise check if field exist in its parent
			parentField, parent_ok := req[field.FieldParent.String].(map[string]interface{})
			ok = parent_ok
			// if field exist check if it is not empty
			if ok {
				_, ok = parentField[field.Field.String]
				if ok {
					ok = (parentField[field.Field.String] != "" && parentField[field.Field.String] != "null" && parentField[field.Field.String] != nil)
				}
			}
		}

		if field.Required.Bool && !ok {
			if valid {
				errRes.GoCode = 35
				valid = false
			} else {
				validationMsg += ", "
			}
			validationMsg += fmt.Sprintf("%s is missing", field.Field.String)
		} else if field.Forbidded.Bool && ok {
			if valid {
				errRes.GoCode = 35
				valid = false
			} else {
				validationMsg += ", "
			}
			validationMsg += fmt.Sprintf("%s is forbided", field.Field.String)
		}
	}

	errRes.Err = errors.New(validationMsg)

	if !valid {
		return
	}

	// check if field that not needed exist to prevent json injection
	for reqField, _ := range req {
		valid2 := false
		for _, conField := range configs {
			if reqField == conField.Field.String {
				valid2 = true
				break
			}
		}

		if valid && !valid2 {
			errRes.GoCode = 35
			valid = false
		}

		if !valid2 && validationMsg != "" {
			validationMsg += ", "
		}

		if !valid2 && validationMsg != "" {
		}

		if !valid2 {
			validationMsg += fmt.Sprintf("%s tidak dibutuhkan", reqField)
		}

	}

	errRes.Err = errors.New(validationMsg)

	return
}
