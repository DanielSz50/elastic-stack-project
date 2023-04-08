package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

type DeleteUploadURI struct {
	uuid string `uri:"uuid"`
}

func registerHandlers(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/upload", func(c *gin.Context) {
		rand.Seed(time.Now().Unix())
		x := rand.Intn(10)
		if x < 3 {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("sample error"))
			return
		}

		c.String(http.StatusOK, "done")
	})

	r.DELETE("/upload/:uuid", func(c *gin.Context) {
		var uri DeleteUploadURI
		if err := c.BindUri(&uri); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if len(uri.uuid) != 2 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.String(http.StatusOK, "upload deleted")
	})
}

func main() {
	const logPath = "/app/logs/gin.log"

	file, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logPath}
	config.EncoderConfig = ecszap.ECSCompatibleEncoderConfig(config.EncoderConfig)
	logger, err := config.Build(ecszap.WrapCoreOption(), zap.AddCaller())
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		httpVersion := fmt.Sprintf("%d.%d", param.Request.ProtoMajor, param.Request.ProtoMinor)

		logLevel := zap.InfoLevel
		if param.ErrorMessage != "" {
			logLevel = zap.ErrorLevel
		}

		logger.Log(logLevel, "",
			zap.String("client.ip", param.ClientIP),
			zap.String("http.request.method", param.Method),
			zap.String("url.path", param.Path),
			zap.String("http.version", httpVersion),
			zap.Int("http.response.status_code", param.StatusCode),
			zap.Int("http.response.body.bytes", param.BodySize),
			zap.String("user_agent.original", param.Request.UserAgent()),
			zap.String("error.message", param.ErrorMessage),
			zap.Int64("event.duration", param.Latency.Nanoseconds()),
			zap.String("event.timezone", "Europe/Warsaw"),
		)

		return ""
	}))

	registerHandlers(r)
	r.Run(":8080")
}
