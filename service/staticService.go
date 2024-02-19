package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/models"
	"alfamart-channel/repo"
	"alfamart-channel/tool"
	"alfamart-channel/util"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zenmuharom/zenlogger"
)

type DefaultStaticService struct {
	logger zenlogger.Zenlogger
	trxLog *domain.Trx
}

type StaticService interface {
	Inquiry(request models.InquiryReq) (response string, err error)
	Payment(request models.PaymentReq) (response string, err error)
	Commit(request models.CommitReq) (response string, err error)
}

func NewStaticService(logger zenlogger.Zenlogger, trxLog *domain.Trx) StaticService {
	return &DefaultStaticService{
		logger: logger,
		trxLog: trxLog,
	}
}

func (this *DefaultStaticService) writeLog() (errRes error) {
	this.logger.Debug("writeLog", zenlogger.ZenField{Key: "trxLog", Value: this.trxLog})
	trxLog := tool.NewTrxLog(this.logger)
	if err := trxLog.Write(this.trxLog); err != nil {
		errDesc := this.logger.Error("writeLog", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call if err := trxLog.Write(\"\", \"\"); err != nil {"})
		errRes = errors.New(errDesc)
		return
	}

	return
}

func (service *DefaultStaticService) mappingRC(productCode, goCode string) (resultCode, resultDesc string) {

	rcConfigRepo := repo.NewRcConfigRepo(service.logger)
	rcConfig, err := rcConfigRepo.FindRc(productCode, goCode)
	if err != nil {
		service.logger.Error("Inquiry", zenlogger.ZenField{Key: "error", Value: err.Error()})
		resultCode = "1706"
		resultDesc = err.Error()
		return
	}

	resultCode = rcConfig.Code.String
	resultDesc = rcConfig.DescEng.String

	return
}

func (service *DefaultStaticService) Inquiry(request models.InquiryReq) (response string, err error) {

	userProductRepo := repo.NewUserProductRepo(service.logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: request.AgentID}, ProductCode: sql.NullString{String: request.ProductID}})

	timeStamp, _ := time.Parse("20060102150405", request.DateTimeRequest)

	service.trxLog.SourceMerchant = sql.NullString{String: userProductConf.Bit33.String, Valid: true}
	service.trxLog.TargetProduct = sql.NullString{String: userProductConf.ProductCode.String, Valid: true}
	service.trxLog.Amount = sql.NullFloat64{Float64: 0, Valid: true}

	req := map[string]interface{}{
		"userName":         request.AgentID,
		"caCode":           userProductConf.Bit32.String,
		"subcaCode":        userProductConf.Bit33.String,
		"productCode":      userProductConf.ProductCodeMapped.String,
		"channel":          userProductConf.Bit18.String,
		"terminal":         "B009TRKK",
		"terminalName":     "12 504 1398",
		"terminalLocation": "Bidakara Pancoran",
		"transactionType":  "38",
		"billNumber":       request.CustomerID,
		"amount":           "0",
		"feeAmount":        "2500",
		"bit61":            request.CustomerID,
		"traxId":           request.AgentTrxID,
		"timeStamp":        timeStamp.Format("2006-01-02 15:04:05:000"),
		"additions":        request,
		"Signature":        request.Signature,
	}

	tsRepo := repo.NewTSRepo(service.logger)
	tsRes, err := tsRepo.Request2(userProductConf.SwitchingUrl.String, int(userProductConf.TimeoutBiller.Int32), req)
	if err != nil {
		service.logger.Error("Inquiry", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	// mapping RC
	resultCode, resultDesc := service.mappingRC(userProductConf.ProductCodeMapped.String, strings.Trim(fmt.Sprintf("%v", tsRes["resultCode"]), " "))

	targetNumber := userProductConf.Bit33.String + userProductConf.ProductCodeMapped.String + request.CustomerID + "|" + request.AgentID + "#" + request.AgentPIN + "#" + request.AgentStoreID

	service.trxLog.TargetNumber = sql.NullString{String: targetNumber, Valid: true}
	service.trxLog.Bit61 = sql.NullString{String: fmt.Sprintf("%v", tsRes["bit61"]), Valid: true}
	service.trxLog.Rc = sql.NullString{String: resultCode, Valid: true}
	service.trxLog.RcDesc = sql.NullString{String: resultDesc, Valid: true}

	arrRes := []string{}
	if tool.CheckRCStatus(service.logger, resultCode, userProductConf.RCSuccess) {
		customerInformation := fmt.Sprintf(
			"%v#%v#%v",
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[20:45], " "),   // nama pt
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[130:144], " "), // no polisi
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[45:75], " "),   // alamat
		)

		arrRes = []string{
			request.AgentID,      // AgentID
			request.AgentPIN,     // AgentPIN
			request.AgentTrxID,   // AgentTrxID
			request.AgentStoreID, // AgentStoreID
			strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[0:20], " "), // CustomerID
			request.DateTimeRequest,             // DateTimeRequest
			resultCode,                          // resultCode
			resultDesc,                          // resultDesc
			time.Now().Format("20060102150405"), // DatetimeResponse
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[144:147], " "), // PaymentPeriod
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[85:130], " "),  // CustomerName
			customerInformation, // customerInformation
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[147:157], " "),                   // Tgl Jatuh Tempo
			strings.TrimSpace(strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[157:169], "0")), // Amount
			strings.TrimSpace(strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[169:181], "0")), // Charge / Nilai Denda
			strings.TrimSpace(strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[207:219], "0")), // Total / min pembayaran
			strings.TrimSpace(strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[231:243], "0")), // AdminFee / biaya admin
			request.ProductID,
			"1",
		}

		targetNumber = userProductConf.Bit33.String + userProductConf.ProductCodeMapped.String + request.CustomerID + "|" + request.AgentID + "#" + request.AgentPIN + "#" + request.AgentStoreID + "#" + strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[144:147], " ")

		service.trxLog.Status = sql.NullString{String: "approve", Valid: true}

	} else {
		arrRes = []string{
			request.AgentID,                     // AgentID
			request.AgentPIN,                    // AgentPIN
			request.AgentTrxID,                  // AgentTrxID
			request.AgentStoreID,                // AgentStoreID
			request.CustomerID,                  // CustomerID
			request.DateTimeRequest,             // DateTimeRequest
			resultCode,                          // resultCode
			resultDesc,                          // resultDesc
			time.Now().Format("20060102150405"), // DatetimeResponse
			"",                                  // PaymentPeriod
			"",                                  // Amount
			"",                                  // Charge / Nilai Denda
			"",                                  // Total / min pembayaran
			"",                                  // AdminFee / biaya admin
			"",                                  // CustomerName
			"",                                  // customerInformation
			"",                                  // tgl jatuh tempo
			"",                                  // amount / min pembayaran
			request.ProductID,
			"1",
		}
	}

	service.writeLog()

	response = strings.Join(arrRes, "|")

	return
}

func (service *DefaultStaticService) Payment(request models.PaymentReq) (response string, err error) {
	rcProcess := "05"
	userProductRepo := repo.NewUserProductRepo(service.logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: request.AgentID}, ProductCode: sql.NullString{String: request.ProductID}})

	service.trxLog.SourceMerchant = sql.NullString{String: userProductConf.Bit33.String, Valid: true}
	service.trxLog.TargetProduct = sql.NullString{String: userProductConf.ProductCode.String, Valid: true}
	service.trxLog.Amount = sql.NullFloat64{Float64: 0, Valid: true}

	targetNumber := userProductConf.Bit33.String + userProductConf.ProductCodeMapped.String + request.CustomerID + "|" + request.AgentID + "#" + request.AgentPIN + "#" + request.AgentStoreID + "#" + request.PaymentPeriod

	trxRepo := repo.NewTrxRepo(service.logger, util.GetDB())
	trx, err := trxRepo.FindByTargetNumber(targetNumber)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			rcProcess = "05"
			service.logger.Info("Payment", zenlogger.ZenField{Key: "error", Value: err.Error()})
		} else {
			service.logger.Error("Payment", zenlogger.ZenField{Key: "error", Value: err.Error()})
		}
	}

	if err == nil {
		if request.Total == strings.TrimSpace(strings.TrimLeft(trx.Bit61.String[157:169], "0")) {
			rcProcess = "00"
		} else {
			rcProcess = "7050"
		}
	}

	// mapping RC
	resultCode, resultDesc := service.mappingRC(userProductConf.ProductCodeMapped.String, rcProcess)

	targetNumber = userProductConf.Bit33.String + userProductConf.ProductCodeMapped.String + request.CustomerID + "|" + request.AgentID + "#" + request.AgentPIN + "#" + request.AgentStoreID + "#" + request.AgentTrxID
	service.trxLog.TargetNumber = sql.NullString{String: targetNumber, Valid: true}
	service.trxLog.Bit61 = sql.NullString{String: trx.Bit61.String, Valid: true}
	service.trxLog.Rc = sql.NullString{String: resultCode, Valid: true}
	service.trxLog.RcDesc = sql.NullString{String: resultDesc, Valid: true}

	arrRes := []string{}
	if tool.CheckRCStatus(service.logger, resultCode, userProductConf.RCSuccess) {
		arrRes = []string{
			request.AgentID,      // AgentID
			request.AgentPIN,     // AgentPIN
			request.AgentTrxID,   // AgentTrxID
			request.AgentStoreID, // AgentStoreID
			strings.TrimLeft(trx.Bit61.String[0:20], " "),                       // CustomerID
			request.DateTimeRequest,                                             // DateTimeRequest
			strings.TrimRight(trx.Bit61.String[144:147], " "),                   // PaymentPeriod
			strings.TrimSpace(strings.TrimLeft(trx.Bit61.String[157:169], "0")), // Amount
			strings.TrimSpace(strings.TrimLeft(trx.Bit61.String[169:181], "0")), // Charge / Nilai Denda
			strings.TrimSpace(strings.TrimLeft(trx.Bit61.String[207:219], "0")), // Total / min pembayaran
			strings.TrimSpace(strings.TrimLeft(trx.Bit61.String[231:243], "0")), // AdminFee / biaya admin
			resultCode,                          // resultCode
			resultDesc,                          // resultDesc
			time.Now().Format("20060102150405"), // DateTimeResponse
			strings.TrimLeft(trx.Bit61.String[75:85], " "), // Pengesahan
			"Untuk info lebih lanjut buka www.partner.id",  // AdditionalData
			request.ProductID, // ProductID
		}

		service.trxLog.Status = sql.NullString{String: "approve", Valid: true}

	} else {
		arrRes = []string{
			request.AgentID,                     // AgentID
			request.AgentPIN,                    // AgentPIN
			request.AgentTrxID,                  // AgentTrxID
			request.AgentStoreID,                // AgentStoreID
			request.CustomerID,                  // CustomerID
			request.DateTimeRequest,             // DateTimeRequest
			request.PaymentPeriod,               // PaymentPeriod
			request.Amount,                      // Amount
			request.Charge,                      // Charge / Nilai Denda
			request.Total,                       // Total / min pembayaran
			resultCode,                          // resultCode
			resultDesc,                          // resultDesc
			time.Now().Format("20060102150405"), // DateTimeResponse
			"",                                  // Pengesahan
			"",                                  // AdditionalData
			request.ProductID,                   // ProductID
		}
	}

	service.writeLog()

	response = strings.Join(arrRes, "|")

	return
}

func (service *DefaultStaticService) Commit(request models.CommitReq) (response string, err error) {
	rcProcess := "05"
	userProductRepo := repo.NewUserProductRepo(service.logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: request.AgentID}, ProductCode: sql.NullString{String: request.ProductID}})

	timeStamp, _ := time.Parse("20060102150405", request.DateTimeRequest)

	service.trxLog.SourceMerchant = sql.NullString{String: userProductConf.Bit33.String, Valid: true}
	service.trxLog.TargetProduct = sql.NullString{String: userProductConf.ProductCode.String, Valid: true}
	service.trxLog.Amount = sql.NullFloat64{Float64: 0, Valid: true}

	targetNumber := userProductConf.Bit33.String + userProductConf.ProductCodeMapped.String + request.CustomerID + "|" + request.AgentID + "#" + request.AgentPIN + "#" + request.AgentStoreID + "#" + request.AgentTrxID

	trxRepo := repo.NewTrxRepo(service.logger, util.GetDB())
	trx, err := trxRepo.FindByCommit(targetNumber)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			rcProcess = "05"
			service.logger.Info("Commit", zenlogger.ZenField{Key: "error", Value: err.Error()})
		} else {
			service.logger.Error("Commit", zenlogger.ZenField{Key: "error", Value: err.Error()})
		}

		// mapping RC
		resultCode, resultDesc := service.mappingRC(userProductConf.ProductCodeMapped.String, rcProcess)

		arrRes := []string{
			request.AgentID,                     // AgentID
			request.AgentPIN,                    // AgentPIN
			request.AgentTrxID,                  // AgentTrxID
			request.AgentStoreID,                // AgentStoreID / No NPP
			request.CustomerID,                  // CustomerID
			request.DateTimeRequest,             // DateTimeRequest
			resultCode,                          // resultCode
			resultDesc,                          // resultDesc
			time.Now().Format("20060102150405"), // DateTimeResponse
			"",                                  // No Pengesahan
			"",
			request.ProductID,
		}

		service.trxLog.TargetNumber = sql.NullString{String: targetNumber, Valid: true}

		service.writeLog()

		response = strings.Join(arrRes, "|")

		return

	}

	req := map[string]interface{}{
		"userName":         request.AgentID,
		"caCode":           userProductConf.Bit32.String,
		"subcaCode":        userProductConf.Bit33.String,
		"productCode":      userProductConf.ProductCodeMapped.String,
		"channel":          userProductConf.Bit18.String,
		"terminal":         "B009TRKK",
		"terminalName":     "12 504 1398",
		"terminalLocation": "Bidakara Pancoran",
		"transactionType":  "50",
		"billNumber":       request.CustomerID,
		"amount":           strings.TrimSpace(strings.TrimLeft(trx.Bit61.String[157:169], "0")),
		"feeAmount":        "0",
		"bit61":            trx.Bit61.String,
		"traxId":           request.AgentTrxID,
		"timeStamp":        timeStamp.Format("2006-01-02 15:04:05:000"),
		"additions":        request,
		"Signature":        request.Signature,
	}

	tsRepo := repo.NewTSRepo(service.logger)
	tsRes, err := tsRepo.Request2(userProductConf.SwitchingUrl.String, int(userProductConf.TimeoutBiller.Int32), req)
	if err != nil {
		service.logger.Error("Inquiry", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	rcProcess = strings.TrimSpace(fmt.Sprintf("%v", tsRes["resultCode"]))

	// mapping RC
	resultCode, resultDesc := service.mappingRC(userProductConf.ProductCodeMapped.String, rcProcess)

	service.trxLog.TargetNumber = sql.NullString{String: targetNumber, Valid: true}
	service.trxLog.Bit61 = sql.NullString{String: fmt.Sprintf("%v", tsRes["bit61"]), Valid: true}
	service.trxLog.Rc = sql.NullString{String: resultCode, Valid: true}
	service.trxLog.RcDesc = sql.NullString{String: resultDesc, Valid: true}

	arrRes := []string{}
	if tool.CheckRCStatus(service.logger, resultCode, userProductConf.RCSuccess) {
		arrRes = []string{
			request.AgentID,         // AgentID
			request.AgentPIN,        // AgentPIN
			request.AgentTrxID,      // AgentTrxID
			request.AgentStoreID,    // AgentStoreID / No NPP
			request.CustomerID,      // CustomerID
			request.DateTimeRequest, // DateTimeRequest
			resultCode,
			resultDesc,
			time.Now().Format("20060102150405"), // DateTimeResponse
			strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[75:85], " "), // No Pengesahan
			"",
			request.ProductID,
		}

		service.trxLog.Status = sql.NullString{String: "approve", Valid: true}

	} else {
		arrRes = []string{
			request.AgentID,                     // AgentID
			request.AgentPIN,                    // AgentPIN
			request.AgentTrxID,                  // AgentTrxID
			request.AgentStoreID,                // AgentStoreID / No NPP
			request.CustomerID,                  // CustomerID
			request.DateTimeRequest,             // DateTimeRequest
			resultCode,                          // resultCode
			resultDesc,                          // resultDesc
			time.Now().Format("20060102150405"), // DateTimeResponse
			strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[95:105], " "), // No Pengesahan
			"",
			request.ProductID,
		}
	}

	service.writeLog()

	response = strings.Join(arrRes, "|")

	return
}
