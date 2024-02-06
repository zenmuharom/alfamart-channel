package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/models"
	"alfamart-channel/repo"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zenmuharom/zenlogger"
)

type DefaultStaticService struct {
	logger  zenlogger.Zenlogger
	request models.Request
}

type StaticService interface {
	Inquiry() (response string, err error)
	Payment() (response string, err error)
	Commit() (response string, err error)
}

func NewStaticService(logger zenlogger.Zenlogger, request models.Request) StaticService {
	return &DefaultStaticService{
		logger:  logger,
		request: request,
	}
}

func (service *DefaultStaticService) Inquiry() (response string, err error) {

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

	customerInformation := fmt.Sprintf(
		"%v#%v#%v",
		strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[20:45], " "),   // nama pt
		strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[130:144], " "), // no polisi
		strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[45:75], " "),   // alamat
	)

	arrRes := []string{
		service.request.AgentID,                                            // AgentID
		service.request.AgentPIN,                                           // AgentPIN
		service.request.AgentTrxID,                                         // AgentTrxID
		service.request.AgentStoreID,                                       // AgentStoreID
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[0:20], " "),     // CustomerID
		service.request.DateTimeRequest,                                    // DateTimeRequest
		strings.Trim(fmt.Sprintf("%v", tsRes["resultCode"]), " "),          // resultCode
		strings.Trim(fmt.Sprintf("%v", tsRes["resultDesc"]), " "),          // resultDesc
		time.Now().Format("20060102150405"),                                // DatetimeResponse
		strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[144:147], " "), // PaymentPeriod
		strings.TrimRight(fmt.Sprintf("%v", tsRes["bit61"])[85:130], " "),  // CustomerName
		customerInformation,                                                // customerInformation
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[147:157], " "),  // tgl jatuh tempo
		strings.TrimLeft(fmt.Sprintf("%v", tsRes["bit61"])[157:169], " "),  // amount / min pembayaran
	}

	response = strings.Join(arrRes, "|")

	return
}

func (service *DefaultStaticService) Payment() (response string, err error) {

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

func (service *DefaultStaticService) Commit() (response string, err error) {

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
