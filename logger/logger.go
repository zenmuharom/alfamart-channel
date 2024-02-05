package logger

import (
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

func SetupLogger() zenlogger.Zenlogger {
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
	return logger
}
