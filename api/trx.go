package api

import (
	"alfamart-channel/domain"
	"alfamart-channel/dto"
	"alfamart-channel/repo"
	"alfamart-channel/service"
	"alfamart-channel/tool"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

func (server *DefaultServer) Trx(ctx *gin.Context) {

	ctx.String(http.StatusOK, "00|Fake Success")
	return

	endpoint := ctx.Request.URL.RawPath
	logger := server.logger
	logger.WithPid(ctx.Request.Header.Get("pid"))

	// read and unmarshal request body to map
	reqBody, _ := ioutil.ReadAll(ctx.Request.Body)
	var req map[string]interface{}
	if err := json.Unmarshal(reqBody, &req); err != nil {
		sendErrorResponse(ctx, logger, ctx.Request.RequestURI, "", ErrorMsg{GoCode: 7026})
		return
	}

	// get productCode
	productCodeKey, productCodeIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "productCode", req)
	if productCodeKey.Parent == "" && productCodeKey.Field == "" {
		err := errors.New("product code route not set yet by Finnet.")
		sendErrorResponse(ctx, logger, ctx.Request.RequestURI, "", ErrorMsg{GoCode: 1994, Err: err})
		return
	}
	productCode := fmt.Sprintf("%v", productCodeIValue)

	// get user config
	userRepo := repo.NewUserRepo(logger)
	_, userNameIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "userName", req)
	userName := fmt.Sprintf("%v", userNameIValue)
	userConf, err := userRepo.Find(userName)
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

	// use config bit18 if empty
	channelKey, channelIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "channel", req)
	if channelIValue == nil {
		if channelKey.Parent == "" {
			req[channelKey.Field] = userConf.ChannelCode.String
		} else {
			// if it has parent key
			parentRCObject := req[channelKey.Parent]
			valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
			newParentRCObject := make(map[string]interface{})
			switch valueOfVariableParentRC.Kind() {
			case reflect.Map:
				iter := valueOfVariableParentRC.MapRange()
				for iter.Next() {
					if iter.Key().String() == channelKey.Field {
						newParentRCObject[iter.Key().String()] = userConf.ChannelCode.String
					} else {
						newParentRCObject[iter.Key().String()] = iter.Value().Interface()
					}
				}
			}
		}
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

	// find transactionType
	transactionTypeKey, transactionTypeIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "transactionType", req)
	if transactionTypeKey.Parent == "" && transactionTypeKey.Field == "" {
		err := errors.New("Transaction type route not set yet.")
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
	}
	transactionType := fmt.Sprintf("%v", transactionTypeIValue)

	// find amount
	amountKey, amountIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "amount", req)
	amount := fmt.Sprintf("%v", amountIValue)
	if amountKey.Parent == "" && amountKey.Field == "" {
		err := errors.New("amount route not set yet.")
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
	}

	// find billNumber
	billNumberKey, billNumberIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "billNumber", req)
	billNumber := fmt.Sprintf("%v", billNumberIValue)
	if billNumberKey.Parent == "" && billNumberKey.Field == "" {
		err := errors.New("bill number route not set yet.")
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
	}

	// find traxId
	traxIdKey, traxIdIValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "traxId", req)
	if traxIdKey.Parent == "" && traxIdKey.Field == "" {
		err := errors.New("transaction id route not set yet.")
		logger.Debug(err.Error())
		sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
	}
	traxId := fmt.Sprintf("%v", traxIdIValue)

	amountToDebet := 0

	// if transactionType is payment and deposit type then do debet to neva
	var nevaService service.NevaService
	if transactionType == "50" && userConf.IsDeposit.Bool {

		// TODO
		// query fee amount config from DB
		amountInteger, _ := strconv.Atoi(amount)

		// TODO here do query
		feeRepo := repo.NewFeeRepo(logger)
		feeConf, err := feeRepo.Find(userName, productCode, amount)
		if err != nil {
			sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 66, Err: nil}) // TODO
			return
		}

		amountToDebet = tool.CalculateFee(logger, *feeConf, amountInteger)

		debetReq := dto.NevaDebet{
			Username:  userConf.NevaUsername.String,
			Password:  userConf.NevaPassword.String,
			DestAcc:   userProductConf.AccountNumber.String,
			Amount:    strconv.Itoa(amountToDebet),
			PhoneNo:   userConf.AccountNumber.String,
			SourceAcc: userConf.AccountNumber.String,
			MitraCo:   userConf.Mitraco.String,
			ProdCode:  productCode,
			BillingNo: billNumber,
			TraxId:    traxId,
		}

		// process to Neva
		logger.Info("Process to NEVA", zenlogger.ZenField{Key: "debetReq", Value: debetReq})
		nevaService = service.NewNevaService(logger)
		nevaStatus, err := nevaService.Debet(debetReq)
		if err != nil {
			logger.Error(err.Error())
			sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 5, Err: err}) // TODO
			return
		}

		if !nevaStatus {
			logger.Info("Neva Status", zenlogger.ZenField{Key: "status", Value: nevaStatus})
			sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 5, Err: err}) // TODO
			return
		}
	}

	// ======================== BEGIN process to TS Adapter
	assignService := service.NewAssignService(logger)
	coreService := service.NewCoreService(logger)

	process := ""
	if transactionType == "38" {
		process = "inquiry"

		constructedBit61 := ""

		// find bit61 then constructed
		bit61Key, bit61IValue := tool.FindFieldAs(logger, domain.SERVER_REQUEST, endpoint, "bit61", req)
		if bit61IValue != nil {

			// construct it
			constructedBit61 = coreService.ConstructBit61(productCode, fmt.Sprintf("%v", bit61IValue))

			if bit61Key.Parent == "" {
				req[bit61Key.Field] = constructedBit61
			} else {
				// if it has parent key
				parentRCObject := req[bit61Key.Parent]
				valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
				newParentRCObject := make(map[string]interface{})
				switch valueOfVariableParentRC.Kind() {
				case reflect.Map:
					iter := valueOfVariableParentRC.MapRange()
					for iter.Next() {
						if iter.Key().String() == bit61Key.Field {
							newParentRCObject[iter.Key().String()] = constructedBit61
						} else {
							newParentRCObject[iter.Key().String()] = iter.Value().Interface()
						}
					}
				}
			}
		}

		// set constructed billNumber
		if billNumberIValue != nil {

			// check if bit61 not constructed yet
			if constructedBit61 == "" {
				// construct it
				constructedBit61 = coreService.ConstructBit61(productCode, fmt.Sprintf("%v", bit61IValue))
			}

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

	} else if transactionType == "50" {
		process = "payment"
	}

	// assign server Request to TS request
	tsReq, err := assignService.AssignMiddlewareRequest("ts_adapter", process, productCode, req)
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

	logger.Debug("coreResponse", zenlogger.ZenField{Key: "response", Value: coreResponse})

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

	// rc = "99" // FOR TESTING

	// if fmt.Sprintf("%s", coreResponse["resultCode"]) == "00" || fmt.Sprintf("%s", coreResponse["resultCode"]) == "0" {
	if (rcIValue != nil || rcIValue != "") && tool.CheckRCStatus(logger, rc, userProductConf.RCSuccess) {

		// assign value
		assignService := service.NewAssignService(logger)
		serverResponse, err := assignService.AssignServerResponse(productCode, ctx.Request.RequestURI, middlewResponseIDs, coreResponse)
		if err != nil {
			logger.Error(err.Error())
			sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 5, Err: err}) // TODO
			return
		}

		// route produback ct_code mapping from TS adapter to channel
		if userProductConf.ProductCodeMapped.String != "" {
			// find productCode key
			productCodeKey, _ := tool.FindFieldAs(logger, domain.SERVER_RESPONSE, endpoint, "productCode", serverResponse)
			if productCodeKey.Parent == "" {
				serverResponse[productCodeKey.Field] = userProductConf.ProductCodeMapped.String
			} else {
				// if it has parent key
				parentRCObject := serverResponse[productCodeKey.Parent]
				valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
				newParentRCObject := make(map[string]interface{})
				switch valueOfVariableParentRC.Kind() {
				case reflect.Map:
					iter := valueOfVariableParentRC.MapRange()
					for iter.Next() {
						if iter.Key().String() == productCodeKey.Field {
							newParentRCObject[iter.Key().String()] = userProductConf.ProductCodeMapped.String
						} else {
							newParentRCObject[iter.Key().String()] = iter.Value().Interface()
						}
					}
				}
			}
		}

		if transactionType == "38" {
			// find amount key
			_, amountResIValue := tool.FindFieldAs(logger, domain.MIDDLEWARE_RESPONSE, fmt.Sprintf("ts_adapter|%v", process), "amount", coreResponse)
			if amountKey.Parent == "" && amountKey.Field == "" {
				err := errors.New("Middleware amount route not set yet.")
				logger.Debug(err.Error())
				sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 1706, Err: err})
				return
			}

			if amountResIValue == 0 || amountResIValue == "0" {
				sendErrorResponseFinish(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 88}, middlewResponseIDs, coreResponse)
				return
			}
		}

		ctx.JSON(http.StatusOK, serverResponse)
		return
	} else {

		// if transactionType is payment and deposit type then do reversal to neva
		if transactionType == "50" && userConf.IsDeposit.Bool && !tool.CheckRCStatus(logger, rc, userProductConf.RCErrorNonDebit) {

			// process to Neva
			nevaStatus, err := nevaService.Reversal()

			// // FOR DEBUG
			// err = errors.New("fail to do reversal")

			if err != nil {
				logger.Error(err.Error())
				sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 7066, Err: err})
				return
			}

			if !nevaStatus {
				logger.Info("Neva Status", zenlogger.ZenField{Key: "status", Value: nevaStatus})
				sendErrorResponse(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: 7066, Err: err})
				return
			}

		}

		rcInt, _ := strconv.ParseInt(rc, 10, 64)
		sendErrorResponseFinish(ctx, logger, endpoint, productCode, ErrorMsg{GoCode: rcInt, Err: err}, middlewResponseIDs, coreResponse)
		return
	}

}
