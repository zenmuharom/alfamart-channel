package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/models"
	"alfamart-channel/repo"
	"alfamart-channel/tool"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zenmuharom/zenlogger"
)

type DefaultStaticService struct {
	logger  zenlogger.Zenlogger
	request models.Request
	trxLog  *domain.Trx
}

type StaticService interface {
	Inquiry(request models.Request) (response string, err error)
	Payment(request models.Request) (response string, err error)
	Commit(request models.Request) (response string, err error)
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

func (service *DefaultStaticService) Inquiry(request models.Request) (response string, err error) {

	service.request = request

	userProductRepo := repo.NewUserProductRepo(service.logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: service.request.AgentID}, ProductCode: sql.NullString{String: service.request.ProductID}})

	timeStamp, _ := time.Parse("20060102150405", service.request.DateTimeRequest)

	service.trxLog.SourceMerchant = sql.NullString{String: userProductConf.Bit33.String, Valid: true}
	service.trxLog.TargetProduct = sql.NullString{String: userProductConf.ProductCode.String, Valid: true}
	service.trxLog.Amount = sql.NullFloat64{Float64: 0, Valid: true}

	req := map[string]interface{}{
		"userName":         service.request.AgentID,
		"caCode":           userProductConf.Bit32.String,
		"subcaCode":        userProductConf.Bit33.String,
		"productCode":      userProductConf.ProductCodeMapped.String,
		"channel":          userProductConf.Bit18.String,
		"terminal":         "B009TRKK",
		"terminalName":     "12 504 1398",
		"terminalLocation": "Bidakara Pancoran",
		"transactionType":  "38",
		"billNumber":       service.request.CustomerID,
		"amount":           "0",
		"feeAmount":        "2500",
		"bit61":            service.request.CustomerID,
		"traxId":           service.request.AgentTrxID,
		"timeStamp":        timeStamp.Format("2006-01-02 15:04:05:000"),
		"additions":        service.request,
		"Signature":        service.request.Signature,
	}

	tsRepo := repo.NewTSRepo(service.logger)
	tsRes, err := tsRepo.Request2(userProductConf.SwitchingUrl.String, int(userProductConf.TimeoutBiller.Int32), req)
	if err != nil {
		service.logger.Error("Inquiry", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	// TODO mapping RC
	resultCode, resultDesc := service.mappingRC(userProductConf.ProductCodeMapped.String, strings.Trim(fmt.Sprintf("%v", tsRes["resultCode"]), " "))

	targetNumber := userProductConf.Bit33.String + userProductConf.ProductCodeMapped.String + service.request.CustomerID + strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[144:147], " ")
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
			service.request.AgentID,                                        // AgentID
			service.request.AgentPIN,                                       // AgentPIN
			service.request.AgentTrxID,                                     // AgentTrxID
			service.request.AgentStoreID,                                   // AgentStoreID
			strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[0:20], " "), // CustomerID
			service.request.DateTimeRequest,                                // DateTimeRequest
			resultCode,                                                     // resultCode
			resultDesc,                                                     // resultDesc
			time.Now().Format("20060102150405"),                            // DatetimeResponse
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[144:147], " "), // PaymentPeriod
			strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[85:130], " "),  // CustomerName
			customerInformation, // customerInformation
			strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[147:157], " "), // tgl jatuh tempo
			strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[157:169], " "), // amount / min pembayaran
		}

		service.trxLog.Status = sql.NullString{String: "approve", Valid: true}

	} else {
		arrRes = []string{
			service.request.AgentID,             // AgentID
			service.request.AgentPIN,            // AgentPIN
			service.request.AgentTrxID,          // AgentTrxID
			service.request.AgentStoreID,        // AgentStoreID
			service.request.CustomerID,          // CustomerID
			service.request.DateTimeRequest,     // DateTimeRequest
			resultCode,                          // resultCode
			resultDesc,                          // resultDesc
			time.Now().Format("20060102150405"), // DatetimeResponse
			"",                                  // PaymentPeriod
			"",                                  // CustomerName
			"",                                  // customerInformation
			"",                                  // tgl jatuh tempo
			"",                                  // amount / min pembayaran
		}
	}

	service.writeLog()

	response = strings.Join(arrRes, "|")

	return
}

func (service *DefaultStaticService) Payment(request models.Request) (response string, err error) {

	service.request = request

	userProductRepo := repo.NewUserProductRepo(service.logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: service.request.AgentID}, ProductCode: sql.NullString{String: service.request.ProductID}})

	timeStamp, _ := time.Parse("20060102150405", service.request.DateTimeRequest)

	service.trxLog.SourceMerchant = sql.NullString{String: userProductConf.Bit33.String, Valid: true}
	service.trxLog.TargetProduct = sql.NullString{String: userProductConf.ProductCode.String, Valid: true}
	service.trxLog.Amount = sql.NullFloat64{Float64: 0, Valid: true}

	req := map[string]interface{}{
		"userName":         service.request.AgentID,
		"caCode":           userProductConf.Bit32.String,
		"subcaCode":        userProductConf.Bit33.String,
		"productCode":      userProductConf.ProductCodeMapped.String,
		"channel":          userProductConf.Bit18.String,
		"terminal":         "B009TRKK",
		"terminalName":     "12 504 1398",
		"terminalLocation": "Bidakara Pancoran",
		"transactionType":  "38",
		"billNumber":       service.request.CustomerID,
		"amount":           "0",
		"feeAmount":        "2500",
		"bit61":            service.request.CustomerID,
		"traxId":           service.request.AgentTrxID,
		"timeStamp":        timeStamp.Format("2006-01-02 15:04:05:000"),
		"additions":        service.request,
		"Signature":        service.request.Signature,
	}

	tsRepo := repo.NewTSRepo(service.logger)
	tsRes, err := tsRepo.Request2(userProductConf.SwitchingUrl.String, int(userProductConf.TimeoutBiller.Int32), req)
	if err != nil {
		service.logger.Error("Inquiry", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	arrRes := []string{
		service.request.AgentID,                                            // AgentID
		service.request.AgentPIN,                                           // AgentPIN
		service.request.AgentTrxID,                                         // AgentTrxID
		service.request.AgentStoreID,                                       // AgentStoreID
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[0:20], " "),     // CustomerID
		service.request.DateTimeRequest,                                    // DateTimeRequest
		strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[144:147], " "), // PaymentPeriod
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[157:169], "0"),  // Amount
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[169:181], "0"),  // Charge / Nilai Denda
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[207:2019], "0"), // Total / min pembayaran
		fmt.Sprintf("%v", tsRes["resultCode"]),                             // resultCode
		fmt.Sprintf("%v", tsRes["resultDesc"]),                             // resultDesc
		time.Now().Format("20060102150405"),                                // DateTimeResponse
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[75:85], " "),    // Pengesahan
		"Untuk info lebih lanjut buka www.partner.id",                      // AdditionalData
		service.request.ProductID,                                          // ProductID
	}

	response = strings.Join(arrRes, "|")

	return
}

func (service *DefaultStaticService) Commit(request models.Request) (response string, err error) {

	service.request = request

	userProductRepo := repo.NewUserProductRepo(service.logger)
	userProductConf, err := userProductRepo.Find(domain.UserProduct{Username: sql.NullString{String: service.request.AgentID}, ProductCode: sql.NullString{String: service.request.ProductID}})

	timeStamp, _ := time.Parse("20060102150405", service.request.DateTimeRequest)

	req := map[string]interface{}{
		"userName":         service.request.AgentID,
		"caCode":           userProductConf.Bit32.String,
		"subcaCode":        userProductConf.Bit33.String,
		"productCode":      userProductConf.ProductCodeMapped.String,
		"channel":          userProductConf.Bit18.String,
		"terminal":         "B009TRKK",
		"terminalName":     "12 504 1398",
		"terminalLocation": "Bidakara Pancoran",
		"transactionType":  "50",
		"billNumber":       service.request.CustomerID,
		"amount":           "0",
		"feeAmount":        "2500",
		"bit61":            service.request.CustomerID,
		"traxId":           service.request.AgentTrxID,
		"timeStamp":        timeStamp.Format("2006-01-02 15:04:05:000"),
		"additions":        service.request,
		"Signature":        service.request.Signature,
	}

	tsRepo := repo.NewTSRepo(service.logger)
	tsRes, err := tsRepo.Request2(userProductConf.SwitchingUrl.String, int(userProductConf.TimeoutBiller.Int32), req)
	if err != nil {
		service.logger.Error("Inquiry", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	arrRes := []string{
		service.request.AgentID,    // AgentID
		service.request.AgentPIN,   // AgentPIN
		service.request.AgentTrxID, // AgentTrxID
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[0:20], " "), // AgentStoreID / No NPP
		service.request.DateTimeRequest,                                // DateTimeRequest
		fmt.Sprintf("%v", tsRes["resultCode"]),
		fmt.Sprintf("%v", tsRes["resultDesc"]),
		time.Now().Format("20060102150405"),                              // DateTimeResponse
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[95:105], " "), // No Pengesahan
	}

	response = strings.Join(arrRes, "|")

	return
}
