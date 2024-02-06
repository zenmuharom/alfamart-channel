package api

import (
	"alfamart-channel/domain"
	"alfamart-channel/logger"
	"alfamart-channel/models"
	"alfamart-channel/service"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

func (handler *DefaultHandler) StaticInquiry(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()
	var request = models.InquiryReq{}

	logger := logger.SetupLogger()
	logger.WithPid(ctx.Request.Header.Get("pid"))

	trxLogStr := ctx.Request.Header.Get("trx_log")

	// parse trx_log
	var trxLog domain.Trx
	if err := json.Unmarshal([]byte(trxLogStr), &trxLog); err != nil {
		logger.Error("StaticInquiry",
			zenlogger.ZenField{Key: "error", Value: err.Error()},
			zenlogger.ZenField{Key: "addition", Value: "error while call if err := json.Unmarshal([]byte(trxLogStr), &trxLog); err != nil {"},
			zenlogger.ZenField{Key: "trxLogStr", Value: trxLogStr},
		)
	}

	logger.Info("StaticInquiry", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "query", Value: queryParams})

	err := handler.parseInqRequest(&request, queryParams)
	if err != nil {
		logger.Error("StaticInquiry", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	trxLog.SourceCode = sql.NullString{String: "INQUIRY", Valid: true}

	staticService := service.NewStaticService(logger, &trxLog)
	response, err := staticService.Inquiry(request)
	if err != nil {
		logger.Error("StaticInquiry", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	// encode trx_log back
	trxLogByte, errEnc := json.Marshal(trxLog)
	if errEnc != nil {
		logger.Error("Trx",
			zenlogger.ZenField{Key: "error", Value: errEnc.Error()},
			zenlogger.ZenField{Key: "addition", Value: "error while call trxLogStr, errEnc = json.Marshal(trxLog)"},
			zenlogger.ZenField{Key: "trxLog", Value: trxLog},
		)
	}

	ctx.Header("trx_log", string(trxLogByte))
	ctx.String(http.StatusOK, response)
}

func (handler *DefaultHandler) StaticPayment(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()
	var request = models.PaymentReq{}

	logger := logger.SetupLogger()
	logger.WithPid(ctx.Request.Header.Get("pid"))

	trxLogStr := ctx.Request.Header.Get("trx_log")

	// parse trx_log
	var trxLog domain.Trx
	if err := json.Unmarshal([]byte(trxLogStr), &trxLog); err != nil {
		logger.Error("StaticPayment",
			zenlogger.ZenField{Key: "error", Value: err.Error()},
			zenlogger.ZenField{Key: "addition", Value: "error while call if err := json.Unmarshal([]byte(trxLogStr), &trxLog); err != nil {"},
			zenlogger.ZenField{Key: "trxLogStr", Value: trxLogStr},
		)
	}

	logger.Info("StaticInquiry", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "query", Value: queryParams})

	err := handler.parsePayRequest(&request, queryParams)
	if err != nil {
		logger.Error("StaticPayment", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	staticService := service.NewStaticService(logger, &trxLog)
	response, err := staticService.Payment(request)
	if err != nil {
		logger.Error("StaticPayment", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	ctx.String(http.StatusOK, response)
}

func (handler *DefaultHandler) StaticCommit(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()
	var request = models.CommitReq{}

	logger := logger.SetupLogger()
	logger.WithPid(ctx.Request.Header.Get("pid"))

	trxLogStr := ctx.Request.Header.Get("trx_log")

	// parse trx_log
	var trxLog domain.Trx
	if err := json.Unmarshal([]byte(trxLogStr), &trxLog); err != nil {
		logger.Error("StaticCommit",
			zenlogger.ZenField{Key: "error", Value: err.Error()},
			zenlogger.ZenField{Key: "addition", Value: "error while call if err := json.Unmarshal([]byte(trxLogStr), &trxLog); err != nil {"},
			zenlogger.ZenField{Key: "trxLogStr", Value: trxLogStr},
		)
	}

	logger.Info("StaticCommit", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "query", Value: queryParams})

	err := handler.parseCommitRequest(&request, queryParams)
	if err != nil {
		logger.Error("StaticCommit", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	staticService := service.NewStaticService(logger, &trxLog)
	response, err := staticService.Commit(request)
	if err != nil {
		logger.Error("StaticCommit", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}

	ctx.String(http.StatusOK, response)
}

func (handler *DefaultHandler) parseInqRequest(req *models.InquiryReq, queryParams url.Values) (err error) {

	req.AgentID = queryParams.Get("AgentID")
	req.AgentPIN = queryParams.Get("AgentPIN")
	req.AgentStoreID = queryParams.Get("AgentStoreID")
	req.AgentTrxID = queryParams.Get("AgentTrxID")
	req.CustomerID = queryParams.Get("CustomerID")
	req.DateTimeRequest = queryParams.Get("DateTimeRequest")
	req.ProductID = queryParams.Get("ProductID")
	req.Signature = queryParams.Get("Signature")

	return
}

func (handler *DefaultHandler) parsePayRequest(req *models.PaymentReq, queryParams url.Values) (err error) {

	req.AgentID = queryParams.Get("AgentID")
	req.AgentPIN = queryParams.Get("AgentPIN")
	req.AgentStoreID = queryParams.Get("AgentStoreID")
	req.AgentTrxID = queryParams.Get("AgentTrxID")
	req.CustomerID = queryParams.Get("CustomerID")
	req.ProductID = queryParams.Get("ProductID")
	req.CustomerID = queryParams.Get("CustomerID")
	req.DateTimeRequest = queryParams.Get("DateTimeRequest")
	req.PaymentPeriod = queryParams.Get("PaymentPeriod")
	req.Amount = queryParams.Get("Amount")
	req.Charge = queryParams.Get("Charge")
	req.Total = queryParams.Get("Total")
	req.AdminFee = queryParams.Get("adminFee")
	req.Signature = queryParams.Get("Signature")

	return
}

func (handler *DefaultHandler) parseCommitRequest(req *models.CommitReq, queryParams url.Values) (err error) {

	req.AgentID = queryParams.Get("AgentID")
	req.AgentPIN = queryParams.Get("AgentPIN")
	req.AgentStoreID = queryParams.Get("AgentStoreID")
	req.AgentTrxID = queryParams.Get("AgentTrxID")
	req.CustomerID = queryParams.Get("CustomerID")
	req.DateTimeRequest = queryParams.Get("DateTimeRequest")
	req.ProductID = queryParams.Get("ProductID")
	req.Signature = queryParams.Get("Signature")

	return
}
