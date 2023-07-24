package main

import (
	"alfamart-channel/api"
	"alfamart-channel/util"
	"log"
	"path/filepath"
	"runtime"

	"github.com/zenmuharom/zenlogger"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Basepath   = filepath.Dir(b)
)

func init() {
	logger := zenlogger.NewZenlogger("init")
	logger.Info("getting config...")
	err := util.LoadConfig(".")
	if err != nil {
		logger.Error(err.Error())
		log.Fatalln(err)
	}

	err = util.ConnectDB()
	if err != nil {
		logger.Error(err.Error())
		log.Fatalln(err)
	}
	logger.Info("config loaded")
}

func main() {
	logger := zenlogger.NewZenlogger("main")
	logger.Info("Service start")
	server := api.New()
	if err := server.Start(); err != nil {
		logger.Error(err.Error())
	}
}
