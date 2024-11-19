package util

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog() *zap.Logger {
	// Ensure the log folder exists
	logFolder := "log"
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		err := os.Mkdir(logFolder, os.ModePerm)
		if err != nil {
			panic("Failed to create log folder: " + err.Error())
		}
	}

	// Get current date for log filenames
	date := time.Now().Format("2006-01-02")

	// Create log files for each level
	infoLog, _ := os.OpenFile(fmt.Sprintf("%s/info-%s.log", logFolder, date), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errorLog, _ := os.OpenFile(fmt.Sprintf("%s/error-%s.log", logFolder, date), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	debugLog, _ := os.OpenFile(fmt.Sprintf("%s/debug-%s.log", logFolder, date), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// Create WriteSyncers for each level
	infoWS := zapcore.AddSync(infoLog)
	errorWS := zapcore.AddSync(errorLog)
	debugWS := zapcore.AddSync(debugLog)

	// Set encoder configurations
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create cores for each level
	infoCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), infoWS, zapcore.InfoLevel)
	errorCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), errorWS, zapcore.ErrorLevel)
	debugCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), debugWS, zapcore.DebugLevel)

	// Combine cores
	core := zapcore.NewTee(infoCore, errorCore, debugCore)

	// Create logger
	logger := zap.New(core)
	return logger
}
