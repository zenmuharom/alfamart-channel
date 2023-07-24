package api

import (
	"alfamart-channel/domain"
	"alfamart-channel/repo"
	"alfamart-channel/service"
	"alfamart-channel/tool"
	"alfamart-channel/util"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

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
	router := gin.New()
	router.Use(server.Logger())
	router.GET("/service/alive", func(ctx *gin.Context) {
		resp := map[string]string{"status": "OK"}
		ctx.JSON(http.StatusOK, resp)
	})
	router.POST("/", server.validate(), server.Trx)
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

func sendErrorResponse(ctx *gin.Context, logger zenlogger.Zenlogger, endpoint, productCode string, errorMsg ErrorMsg) {

	logger.Debug("sendErrorResponse", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "errorMsg", Value: errorMsg})
	httpStatus := http.StatusInternalServerError
	resultCode := fmt.Sprintf("%v", errorMsg.GoCode)

	assignServResponseService := service.NewAssignService(logger)
	serverResponse, err := assignServResponseService.AssignServerResponse("", endpoint, []int{}, map[string]interface{}{})
	if err != nil {
		logger.Error(err.Error())
	}

	// find RC Code
	rcConfigRepo := repo.NewRcConfigRepo(logger)
	rcConfig, errRC := rcConfigRepo.FindRc(productCode, errorMsg.GoCode) // TODO
	if errRC != nil {
		errRC = errors.New(fmt.Sprintf("RC Config not found: %v", errRC.Error()))
		logger.Error(errRC.Error())
		errorMsg.GoCode = 7000
		resultCode = fmt.Sprintf("%v", errorMsg.GoCode)
		errorMsg.Err = errRC
	} else {
		httpStatus = int(rcConfig.Httpstatus.Int64)
		errorMsg.GoCode, _ = strconv.ParseInt(rcConfig.Code.String, 10, 64)
		resultCode = rcConfig.Code.String
		if errorMsg.Err != nil {
			errorMsg.Err = errors.New(fmt.Sprintf("%v %v", rcConfig.DescEng.String, errorMsg.Err.Error()))
		} else {
			errorMsg.Err = errors.New(rcConfig.DescEng.String)
		}

	}

	// find resultCode
	keyResultCode, _ := tool.FindFieldAs(logger, domain.SERVER_RESPONSE, ctx.Request.RequestURI, "resultCode", serverResponse)
	if keyResultCode.Field == "" && keyResultCode.Parent == "" {
		err := errors.New("Result code not set yet")
		errorMsg.Err = err
		logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorMsg.Err.Error()})
		return
	}

	if keyResultCode.Parent == "" {
		serverResponse[keyResultCode.Field] = fmt.Sprintf("%v", resultCode)
	} else {
		// if it has parent key
		parentRCObject := serverResponse[keyResultCode.Parent]
		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
		newParentRCObject := make(map[string]interface{})
		switch valueOfVariableParentRC.Kind() {
		case reflect.Map:
			iter := valueOfVariableParentRC.MapRange()
			for iter.Next() {
				if iter.Key().String() == keyResultCode.Field {
					newParentRCObject[iter.Key().String()] = fmt.Sprintf("%v", resultCode)
				} else {
					newParentRCObject[iter.Key().String()] = iter.Value().Interface()
				}
			}
		}
	}

	// find resultDesc
	keyResultDesc, _ := tool.FindFieldAs(logger, domain.SERVER_RESPONSE, ctx.Request.RequestURI, "resultDesc", serverResponse)
	if keyResultDesc.Field == "" && keyResultDesc.Parent == "" {
		err := errors.New("Result desc not set yet")
		errorMsg.Err = err
		logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorMsg)
		return
	}

	if keyResultDesc.Parent == "" {
		serverResponse[keyResultDesc.Field] = errorMsg.Err.Error()
	} else {
		// if it has parent key
		parentRCObject := serverResponse[keyResultDesc.Parent]
		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
		newParentRCObject := make(map[string]interface{})
		switch valueOfVariableParentRC.Kind() {
		case reflect.Map:
			iter := valueOfVariableParentRC.MapRange()
			for iter.Next() {
				if iter.Key().String() == keyResultDesc.Field {
					newParentRCObject[iter.Key().String()] = errorMsg.Err.Error()
				} else {
					newParentRCObject[iter.Key().String()] = iter.Value().Interface()
				}
			}
		}
	}

	ctx.AbortWithStatusJSON(httpStatus, serverResponse)
	return

}

func sendErrorResponseFinish(ctx *gin.Context, logger zenlogger.Zenlogger, endpoint, productCode string, errorMsg ErrorMsg, middlewResponseIDs []int, response map[string]interface{}) {

	logger.Debug("sendErrorResponse", zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "errorMsg", Value: errorMsg})
	httpStatus := http.StatusInternalServerError
	resultCode := fmt.Sprintf("%v", errorMsg.GoCode)

	assignServResponseService := service.NewAssignService(logger)
	serverResponse, err := assignServResponseService.AssignServerResponse("", endpoint, middlewResponseIDs, response)
	if err != nil {
		logger.Error(err.Error())
	}

	// find RC Code
	rcConfigRepo := repo.NewRcConfigRepo(logger)
	rcConfig, errRC := rcConfigRepo.FindRc(productCode, errorMsg.GoCode) // TODO
	if errRC != nil {
		errRC = errors.New(fmt.Sprintf("RC Config not found: %v", errRC.Error()))
		logger.Error(errRC.Error())
		errorMsg.GoCode = 7000
		resultCode = fmt.Sprintf("%v", errorMsg.GoCode)
		errorMsg.Err = errRC
	} else {
		httpStatus = int(rcConfig.Httpstatus.Int64)
		errorMsg.GoCode, _ = strconv.ParseInt(rcConfig.Code.String, 10, 64)
		resultCode = rcConfig.Code.String
		if errorMsg.Err != nil {
			errorMsg.Err = errors.New(fmt.Sprintf("%v %v", rcConfig.DescEng.String, errorMsg.Err.Error()))
		} else {
			errorMsg.Err = errors.New(rcConfig.DescEng.String)
		}

	}

	// find resultCode
	keyResultCode, _ := tool.FindFieldAs(logger, domain.SERVER_RESPONSE, ctx.Request.RequestURI, "resultCode", serverResponse)
	if keyResultCode.Field == "" && keyResultCode.Parent == "" {
		err := errors.New("Result code not set yet")
		errorMsg.Err = err
		logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorMsg.Err.Error()})
		return
	}

	if keyResultCode.Parent == "" {
		serverResponse[keyResultCode.Field] = fmt.Sprintf("%v", resultCode)
	} else {
		// if it has parent key
		parentRCObject := serverResponse[keyResultCode.Parent]
		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
		newParentRCObject := make(map[string]interface{})
		switch valueOfVariableParentRC.Kind() {
		case reflect.Map:
			iter := valueOfVariableParentRC.MapRange()
			for iter.Next() {
				if iter.Key().String() == keyResultCode.Field {
					newParentRCObject[iter.Key().String()] = fmt.Sprintf("%v", resultCode)
				} else {
					newParentRCObject[iter.Key().String()] = iter.Value().Interface()
				}
			}
		}
	}

	// find resultDesc
	keyResultDesc, _ := tool.FindFieldAs(logger, domain.SERVER_RESPONSE, ctx.Request.RequestURI, "resultDesc", serverResponse)
	if keyResultDesc.Field == "" && keyResultDesc.Parent == "" {
		err := errors.New("Result desc not set yet")
		errorMsg.Err = err
		logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorMsg)
		return
	}

	if keyResultDesc.Parent == "" {
		serverResponse[keyResultDesc.Field] = errorMsg.Err.Error()
	} else {
		// if it has parent key
		parentRCObject := serverResponse[keyResultDesc.Parent]
		valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
		newParentRCObject := make(map[string]interface{})
		switch valueOfVariableParentRC.Kind() {
		case reflect.Map:
			iter := valueOfVariableParentRC.MapRange()
			for iter.Next() {
				if iter.Key().String() == keyResultDesc.Field {
					newParentRCObject[iter.Key().String()] = errorMsg.Err.Error()
				} else {
					newParentRCObject[iter.Key().String()] = iter.Value().Interface()
				}
			}
		}
	}

	ctx.AbortWithStatusJSON(httpStatus, serverResponse)
	return

}
