package api

import (
	"alfamart-channel/logger"
	"alfamart-channel/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

func (handler *DefaultHandler) StaticInquiry(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()

	logger := logger.SetupLogger()
	logger.WithPid(ctx.Request.Header.Get("pid"))

	logger.Info("Static", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "query", Value: queryParams})

	err := handler.parseRequest(&handler.request, queryParams)
	if err != nil {
		logger.Error("General", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	staticService := service.NewStaticService(logger, handler.request)
	response, err := staticService.Inquiry()
	if err != nil {

	}

	ctx.String(http.StatusOK, response)
}

func (handler *DefaultHandler) StaticPayment(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()

	logger := logger.SetupLogger()
	logger.WithPid(ctx.Request.Header.Get("pid"))

	logger.Info("Static", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "query", Value: queryParams})

	err := handler.parseRequest(&handler.request, queryParams)
	if err != nil {
		logger.Error("General", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	staticService := service.NewStaticService(logger, handler.request)
	response, err := staticService.Payment()
	if err != nil {

	}

	ctx.String(http.StatusOK, response)
}

func (handler *DefaultHandler) StaticCommit(ctx *gin.Context) {

	endpoint := ctx.Request.URL.Path
	queryParams := ctx.Request.URL.Query()

	logger := logger.SetupLogger()
	logger.WithPid(ctx.Request.Header.Get("pid"))

	logger.Info("Static", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "query", Value: queryParams})

	err := handler.parseRequest(&handler.request, queryParams)
	if err != nil {
		logger.Error("General", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	staticService := service.NewStaticService(logger, handler.request)
	response, err := staticService.Commit()
	if err != nil {

	}

	ctx.String(http.StatusOK, response)
}
