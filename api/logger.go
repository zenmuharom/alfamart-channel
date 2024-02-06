package api

import (
	"alfamart-channel/domain"
	"alfamart-channel/logger"
	"alfamart-channel/repo"
	"alfamart-channel/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zenmuharom/zenlogger"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (server *DefaultServer) Logger() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w

		serveTime := time.Now()

		logger := logger.SetupLogger()

		body, _ := io.ReadAll(ctx.Request.Body)
		fields := []zenlogger.ZenField{
			{
				Key:   "URI",
				Value: ctx.Request.RequestURI,
			},
			{
				Key:   "ClientIP",
				Value: ctx.ClientIP(),
			},
			{
				Key:   "RemoteIP",
				Value: ctx.RemoteIP(),
			},
			{
				Key:   "Header",
				Value: ctx.Request.Header,
			},
			{
				Key:   "Body",
				Value: string(body),
			},
		}
		logger.Info("Request", fields...)

		// write log
		trxLog, err := server.insertTrx(logger)
		if err != nil {
			errDesc := logger.Error("Logger", zenlogger.ZenField{Key: "error", Value: err.Error()})
			errRes := errors.New(errDesc)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errRes)
			return
		}

		trxLogStr, errMarshall := json.Marshal(trxLog)
		if errMarshall != nil {
			errDesc := logger.Error("Logger", zenlogger.ZenField{Key: "error", Value: errMarshall.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while marshal trxLog"})
			errRes := errors.New(errDesc)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errRes)
			return
		}

		// add pid & trx_log
		ctx.Request.Header.Add("pid", logger.GetPid())
		ctx.Request.Header.Add("trx_log", string(trxLogStr))
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// go to handler
		ctx.Next()

		// hook middleware before send response

		fields = []zenlogger.ZenField{
			{
				Key:   "HTTPStatus",
				Value: strconv.Itoa(ctx.Writer.Status()),
			},
			{
				Key:   "Header",
				Value: w.Header(),
			},
			{
				Key:   "Body",
				Value: w.body.String(),
			},
		}
		logger.Info("Response", fields...)

		responseTime := time.Now()
		elapsedTime := time.Since(serveTime)
		fields = []zenlogger.ZenField{
			{
				Key:   "RequestTime",
				Value: serveTime.Format("2006-01-02T15:04:05-0700"),
			},
			{
				Key:   "ResponseTime",
				Value: responseTime.Format("2006-01-02T15:04:05-0700"),
			},
			{
				Key:   "Latency",
				Value: fmt.Sprintf("%sms", strconv.Itoa(int(elapsedTime.Milliseconds()))),
			},
		}

		if err := json.Unmarshal([]byte(w.Header().Get("trx_log")), trxLog); err != nil {
			logger.Error("Logger", zenlogger.ZenField{Key: "error", Value: err.Error()})
		}

		trxLog.ElapsedTime = sql.NullInt64{Int64: elapsedTime.Milliseconds(), Valid: true}

		server.updateTrx(logger, trxLog)

		logger.Info("Summary", fields...)
	}
}

func (server *DefaultServer) insertTrx(logger zenlogger.Zenlogger) (trxLog *domain.Trx, err error) {

	dataToInsert := domain.Trx{
		Pid:            logger.GetPid(),
		AmendmentDate:  sql.NullTime{Time: time.Now(), Valid: true},
		SourceMerchant: sql.NullString{String: "000000", Valid: true},
		TargetProduct:  sql.NullString{String: "000000", Valid: true},
		Status:         sql.NullString{String: "decline", Valid: true},
		ElapsedTime:    sql.NullInt64{Int64: 0, Valid: true},
	}

	trxRepo := repo.NewTrxRepo(logger, util.GetDB())
	trxLog, err = trxRepo.Upsert(dataToInsert)
	if err != nil {
		logger.Error("insertTrx", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call _, err := trxRepo.Upsert(dataToInsert)"})
	}

	return
}

func (server *DefaultServer) updateTrx(logger zenlogger.Zenlogger, trxLog *domain.Trx) (err error) {

	trxRepo := repo.NewTrxRepo(logger, util.GetDB())
	_, err = trxRepo.Upsert(*trxLog)
	if err != nil {
		logger.Error("updateTrx", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call _, err := trxRepo.Upsert(trxLog)"})
	}

	return
}
