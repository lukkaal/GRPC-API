package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// logger for gateway(one of the micro-service)
var GinloggerObj *logrus.Logger

// init when import one package
func init() {
	// new logrus instance
	logger := logrus.New()

	logfile, err := setOutputFile()
	if err != nil {
		// false: stdout
		logger.Out = os.Stdout
		fmt.Println("logger init failed:", err)
	} else {
		logger.Out = logfile
	}

	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	// successfully initiated
	GinloggerObj = logger
}

// set output file for logger instance
func setOutputFile() (*os.File, error) {
	// set path to GRPC-API/logs
	logDir := "logs"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf(
				"failed to create log directory: %v", err)
		}
	}

	currentTime := time.Now().Format("2006-01-02")
	logFileName := path.Join(logDir, currentTime+".log")

	// opwn(or create) file
	file, err := os.OpenFile(
		logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	return file, nil
}
