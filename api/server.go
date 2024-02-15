package api

import (
	"alfamart-channel/models"
	"alfamart-channel/util"
	"fmt"
	"strings"
	"time"

	"net/http"

	"github.com/zenmuharom/zenlogger"

	"github.com/gin-gonic/gin"
)

type ErrorMsg struct {
	GoCode int64
	Err    error
}

type DefaultServer struct {
	router *gin.Engine
	logger zenlogger.Zenlogger
}

type Server interface {
	Start() error
	setupRouter()
	setupLogger() zenlogger.Zenlogger
}

func New() (server Server) {
	gin.SetMode(util.GetConfig().GIN_MODE)
	server = &DefaultServer{
		router: gin.New(),
	}
	server.setupLogger()
	server.setupRouter()
	return
}

func (server *DefaultServer) setupLogger() zenlogger.Zenlogger {
	logger := zenlogger.NewZenlogger()
	config := zenlogger.Config{
		Pid: zenlogger.ZenConf{
			Label: "traceId",
		},
		Severity: zenlogger.Severity{
			Label:  "severity",
			Access: "ACCESS",
			Info:   "INFO",
			Debug:  "DEBUG",
			Error:  "ERROR",
			Query:  "QUERY",
		},
	}

	if util.GetConfig().ENV == "local" {
		config.Output = zenlogger.Output{
			Path:   "logs",
			Format: "2006-01-02 15",
		}
	}

	if util.GetConfig().ENV == "prod" {
		config.Production = true
	}

	logger.SetConfig(config)
	server.logger = logger
	return logger
}

func (server *DefaultServer) setupRouter() {

	handler := NewHandler()

	router := gin.New()
	router.GET("/service/alive", func(ctx *gin.Context) {
		resp := map[string]string{"status": "OK"}
		ctx.JSON(http.StatusOK, resp)
	})

	trx := router.Group("")
	trx.Use(server.Logger())

	inquiry := trx.Group("")
	inquiry.Use(server.validate()).GET("/adira/inquiry", handler.StaticInquiry)

	payment := trx.Group("")
	payment.Use(server.validatePayment()).GET("/adira/payment", handler.StaticPayment)

	commit := trx.Group("")
	commit.Use(server.validate()).GET("/adira/commit", handler.StaticCommit)

	// handlers := map[string]gin.HandlerFunc{
	// 	"general": handler.General,
	// }

	// routeRepo := repo.NewRouteRepo(server.logger)
	// routes, err := routeRepo.FindAll()
	// if err != nil {
	// 	server.logger.Error("there is no route config has been created", zenlogger.ZenField{Key: "error", Value: err.Error()})
	// 	os.Exit(0)
	// }

	// for _, route := range routes {
	// 	switch route.Method.String {
	// 	case "GET":
	// 		router.GET(route.Path.String, server.validate(), handlers[route.Handler.String])
	// 	case "POST":
	// 		router.POST(route.Path.String, server.validate(), handlers[route.Handler.String])
	// 	case "PUT":
	// 		router.PUT(route.Path.String, server.validate(), handlers[route.Handler.String])
	// 	case "DELETE":
	// 		router.DELETE(route.Path.String, server.validate(), handlers[route.Handler.String])
	// 	}
	// 	server.logger.Info("starting endpoint", zenlogger.ZenField{Key: "method", Value: route.Method.String}, zenlogger.ZenField{Key: "path", Value: route.Path.String})
	// }
	server.router = router
}

func (server *DefaultServer) Start() error {
	config := util.GetConfig()
	server.logger.Info(fmt.Sprintf("service started at %s:%s", config.ServerAddress, config.ServerPort))
	server.router.Run(fmt.Sprintf("%s:%s", config.ServerAddress, config.ServerPort))
	return nil
}

type ValidationErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func sendErrorResponse(ctx *gin.Context, logger zenlogger.Zenlogger, endpoint, productCode string, errorMsg models.ErrorMsg) {

	logger.Debug("sendErrorResponse", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "errorMsg", Value: errorMsg})

	arrRes := []string{
		"",                                  // AgentID
		"",                                  // AgentPIN
		"",                                  // AgentTrxID
		"",                                  // AgentStoreID
		"",                                  // CustomerID
		"",                                  // DateTimeRequest
		"05",                                // resultCode
		errorMsg.Err.Error(),                // resultDesc
		time.Now().Format("20060102150405"), // DatetimeResponse
		"",                                  // PaymentPeriod
		"",                                  // CustomerName
		"",                                  // customerInformation
		"",                                  // tgl jatuh tempo
		"",                                  // amount / min pembayaran
	}

	response := strings.Join(arrRes, "|")

	ctx.Header("trx_log", ctx.Request.Header.Get("trx_log"))

	ctx.String(http.StatusOK, response)
	ctx.Abort()
	return

}

func writeToText(json []interface{}) (responseText string) {

	index := 0
	for _, value := range json {
		if value == nil {
			value = ""
		}
		responseText += fmt.Sprintf("%v", value)
		if index < len(json) {
			responseText += "|"
		}
		index++
	}

	return
}
