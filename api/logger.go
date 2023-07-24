package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

		logger := server.setupLogger()

		body, _ := ioutil.ReadAll(ctx.Request.Body)
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
		ctx.Request.Header.Add("pid", logger.GetPid())
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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

		logger.Info("Summary", fields...)
	}
}
