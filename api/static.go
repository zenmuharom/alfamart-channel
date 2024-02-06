package api

import (
	"alfamart-channel/domain"
	"alfamart-channel/logger"
	"alfamart-channel/service"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

func (handler *DefaultHandler) StaticInquiry(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()

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

	err := handler.parseRequest(&handler.request, queryParams)
	if err != nil {
		logger.Error("StaticInquiry", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	trxLog.SourceCode = sql.NullString{String: "INQUIRY", Valid: true}

	staticService := service.NewStaticService(logger, &trxLog)
	response, err := staticService.Inquiry(handler.request)
	if err != nil {

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

	err := handler.parseRequest(&handler.request, queryParams)
	if err != nil {
		logger.Error("StaticPayment", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	staticService := service.NewStaticService(logger, &trxLog)
	response, err := staticService.Payment(handler.request)
	if err != nil {

	}

	ctx.String(http.StatusOK, response)
}

func (handler *DefaultHandler) StaticCommit(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()

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

	err := handler.parseRequest(&handler.request, queryParams)
	if err != nil {
		logger.Error("StaticCommit", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	staticService := service.NewStaticService(logger, &trxLog)
	response, err := staticService.Commit(handler.request)
	if err != nil {

	}

	ctx.String(http.StatusOK, response)
}
