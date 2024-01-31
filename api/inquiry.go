package api

import (
	"alfamart-channel/domain"
	"alfamart-channel/repo"
	"alfamart-channel/service"
	"alfamart-channel/tool"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

func (server *DefaultServer) Inquiry(ctx *gin.Context) {
	// endpoint := ctx.Request.URL.Path
	endpoint := ctx.Request.URL.Path
	logger := server.logger
	logger.WithPid(ctx.Request.Header.Get("pid"))
	process := "inquiry"
	logger.Info(process, zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "queryParams", Value: ctx.Request.URL.Query()})

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

	// get productCode
	productCodeKey, productCodeIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "productCode", req)
	if productCodeKey.Parent == "" && productCodeKey.Field == "" {
		err := errors.New("product code route not set yet by Finnet.")
		sendErrorResponse(ctx, logger, endpoint, "", ErrorMsg{GoCode: 1994, Err: err})
		return
	}
	productCode := fmt.Sprintf("%v", productCodeIValue)

	// get user config
	userRepo := repo.NewUserRepo(logger)
	_, userNameIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "userName", req)
	userName := fmt.Sprintf("%v", userNameIValue)
	_, err := userRepo.Find(userName)
	if err != nil {
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err}) // TODO
		return
	}

	// validating IP based on rules on DB
	ipConfigRepo := repo.NewIpConfigRepo(logger)
	validIP, err := ipConfigRepo.FindIpByIpUsernameAndProductCode(ctx.ClientIP(), userName, productCode)
	if err != nil {
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 5, Err: err}) // TODO
		return
	}

	if !validIP {
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1994, Err: nil}) // TODO
		return
	}

	// get user product config
	userProductRepo := repo.NewUserProductRepo(logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: userName}, ProductCode: sql.NullString{String: productCode}})
	if err != nil {
		logger.Error(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 17, Err: nil})
		return
	}

	if userProductConf.RCSuccess == nil {
		err = errors.New("RC mapping not set yet by Finnet.")
		logger.Error(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
		return
	}

	// route product_code mapping to TS adapter
	if userProductConf.ProductCodeMapped.String != "" {
		productCode = userProductConf.ProductCodeMapped.String
	}

	// find billNumber
	billNumberKey, billNumberIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "billNumber", req)
	billNumber := fmt.Sprintf("%v", billNumberIValue)
	if billNumberKey.Parent == "" && billNumberKey.Field == "" {
		err := errors.New("bill number route not set yet.")
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
	}

	// ======================== BEGIN process to TS Adapter
	assignService := service.NewAssignService(logger, endpoint, productCode)
	coreService := service.NewCoreService(logger)

	constructedBit61 := ""

	// set constructed billNumber
	if billNumberIValue != nil {

		constructedBit61 = coreService.ConstructBit61(productCode, billNumber)

		if billNumberKey.Parent == "" {
			req[billNumberKey.Field] = constructedBit61
		} else {
			// if it has parent key
			parentRCObject := req[billNumberKey.Parent]
			valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
			newParentRCObject := make(map[string]interface{})
			switch valueOfVariableParentRC.Kind() {
			case reflect.Map:
				iter := valueOfVariableParentRC.MapRange()
				for iter.Next() {
					if iter.Key().String() == billNumberKey.Field {
						newParentRCObject[iter.Key().String()] = constructedBit61
					} else {
						newParentRCObject[iter.Key().String()] = iter.Value().Interface()
					}
				}
			}
		}
	}

	// assign server Request to TS request
	tsReq, err := assignService.AssignMiddlewareRequest("ts_adapter", process, &userProductConf, req)
	if err != nil {
		logger.Error(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 7000, Err: err})
		return
	}

	// prepare to receive response
	var coreResponse map[string]interface{}
	middlewResponseIDs, coreResponse, err := coreService.Request2(process, &userProductConf, tsReq)
	if err != nil {
		logger.Error(err.Error())
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 68, Err: nil}) // TODO
			return
		} else if err == sql.ErrNoRows {
			sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err}) // TODO
			return
		}

		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 7000, Err: err}) // TODO
		return
	}
	// ======================== END process to TS Adapter

	// find rc key
	rcKey, rcIValue := tool.FindFieldAs(logger, domain.MIDDLEWARE_RESPONSE, fmt.Sprintf("ts_adapter|%v", process), "rc", coreResponse)
	if rcKey.Parent == "" && rcKey.Field == "" {
		err := errors.New("Middleware RC route not set yet.")
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
		return
	}

	// if rc is null or empty string then set it to 7000: "System Maintenance"
	if rcIValue == nil || rcIValue == "" {
		rcIValue = 7000
	}

	rc := fmt.Sprintf("%v", rcIValue)

	// assign value
	serverResponse, err := assignService.AssignServerResponse(productCode, endpoint, middlewResponseIDs, coreResponse)
	if err != nil {
		logger.Error(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 5, Err: err})
		return
	}

	serverResponseText, err := assignService.OrderingServerResponse(serverResponse)
	if err != nil {
		logger.Error("Inquiry", zenlogger.ZenField{Key: "err", Value: err.Error()})
	}

	if (rcIValue != nil || rcIValue != "") && tool.CheckRCStatus(logger, rc, userProductConf.RCSuccess) {

		fmt.Println(logger.Debug("serverResponse", zenlogger.ZenField{Key: "status", Value: "success"}, zenlogger.ZenField{Key: "value", Value: serverResponse}))

		ctx.String(http.StatusOK, serverResponseText)
	} else {
		ctx.String(http.StatusOK, serverResponseText)
	}

	return

}
