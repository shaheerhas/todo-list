package app

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerToFile() gin.HandlerFunc {
	//logFilePath := os.Getenv("LOGS_PATH")
	logFileName := os.Getenv("LOG_FILENAME")
	//log file
	//	fileName := path.Join(logFilePath, logFileName)
	//fileName := path.Join(logFileName, logFileName)
	//write file
	src, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	//instantiation

	logger := logrus.New()
	//Set output
	logger.Out = src
	//Set log level
	logger.SetLevel(logrus.DebugLevel)
	//Format log
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return func(c *gin.Context) {
		// start time
		startTime := time.Now()
		// Processing request
		c.Next()
		// End time
		endTime := time.Now()
		// execution time
		latencyTime := endTime.Sub(startTime)
		// Request mode
		reqMethod := c.Request.Method
		// Request routing
		reqUri := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Request IP
		clientIP := c.ClientIP()
		// Log format
		if statusCode >= 500 && statusCode < 600 {
			logger.Errorf("| %3d | %13v | %15s | %s | %s |",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri,
			)
		} else if statusCode >= 400 && statusCode < 500 {
			logger.Warnf("| %3d | %13v | %15s | %s | %s |",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri,
			)
		} else if statusCode >= 200 && statusCode < 300 {
			logger.Infof("| %3d | %13v | %15s | %s | %s |",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri,
			)
		}
	}
}
