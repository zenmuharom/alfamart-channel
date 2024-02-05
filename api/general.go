package api

import (
	"alfamart-channel/logger"
	"alfamart-channel/models"
	"alfamart-channel/service"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

type DefaultHandler struct {
	request models.Request
}

type Handler interface {
	General(ctx *gin.Context)
	Static(ctx *gin.Context)
}

func NewHandler() Handler {
	return &DefaultHandler{}
}

func (handler *DefaultHandler) General(ctx *gin.Context) {
	endpoint := ctx.Request.URL.Path
	logger := logger.SetupLogger()
	logger.WithPid(ctx.Request.Header.Get("pid"))

	queryParams := ctx.Request.URL.Query()
	// Iterate through the map and print the key-value pairs
	reqMap := make(map[string]interface{})
	err := handler.parseRequestToRaw(reqMap, queryParams)
	if err != nil {
		logger.Error("General", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	logger.Info("General", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "query", Value: reqMap})

	err = handler.parseRequest(&handler.request, queryParams)
	if err != nil {
		logger.Error("General", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while decode request body to map"})
	}

	routeService := service.NewRoute(logger, endpoint, handler.request, reqMap)
	response, err := routeService.Process()
	logger.Debug("General", zenlogger.ZenField{Key: "response", Value: response})
	if err != nil {
		errDesc := logger.Error("General", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error during the process"})
		ctx.String(http.StatusOK, errDesc)
		return
	}

	orderedResponse, err := routeService.OrderingServerResponse(response)
	if err != nil {
		errDesc := logger.Error("General", zenlogger.ZenField{Key: "error", Value: err.Error()})
		ctx.String(http.StatusOK, errDesc)
	}

	logger.Info("General", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "queryParams", Value: ctx.Request.URL.Query()})

	ctx.String(http.StatusOK, orderedResponse)

}

func (handler *DefaultHandler) parseRequest(req *models.Request, queryParams url.Values) (err error) {

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

func (handler *DefaultHandler) parseRequestToRaw(req map[string]interface{}, queryParams url.Values) (err error) {
	for key, values := range queryParams {
		// If there is only one value, add it directly to the map
		if len(values) == 1 {
			req[key] = values[0]
		} else {
			// If there are multiple values, add them as a slice of strings
			req[key] = values
		}
	}

	return
}
