package util

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(config Configuration) *zap.Logger {
	// Ensure the log root folder exists
	logFolder := config.Dir.Logs
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		err := os.Mkdir(logFolder, os.ModePerm)
		if err != nil {
			panic("Failed to create log folder: " + err.Error())
		}
	}

	// Get current date for log filenames
	date := time.Now().Format("2006-01-02")

	// Create folders for each log level
	logLevels := []string{"info", "error", "debug"}
	for _, level := range logLevels {
		levelFolder := fmt.Sprintf("%s/%s", logFolder, level)
		if _, err := os.Stat(levelFolder); os.IsNotExist(err) {
			err := os.Mkdir(levelFolder, os.ModePerm)
			if err != nil {
				panic("Failed to create log level folder: " + err.Error())
			}
		}
	}

	// Create log files for each level, with filename based on date
	infoLog, _ := os.OpenFile(fmt.Sprintf("%s/info/%s-%s.log", logFolder, "info", date), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errorLog, _ := os.OpenFile(fmt.Sprintf("%s/error/%s-%s.log", logFolder, "error", date), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	debugLog, _ := os.OpenFile(fmt.Sprintf("%s/debug/%s-%s.log", logFolder, "debug", date), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// Create WriteSyncers for each level
	infoWS := zapcore.AddSync(infoLog)
	errorWS := zapcore.AddSync(errorLog)
	debugWS := zapcore.AddSync(debugLog)

	// Add a console output for debug logs if DEBUG is true
	consoleDebugWS := zapcore.AddSync(os.Stdout)

	// Set encoder configurations
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create cores for each level
	infoCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), infoWS, zapcore.InfoLevel)
	errorCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), errorWS, zapcore.ErrorLevel)

	// Use Tee for debug logs to write to both file and console if DEBUG is true
	debugCore := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), debugWS, zapcore.DebugLevel),
	)
	if config.Debug {
		debugCore = zapcore.NewTee(
			debugCore,
			zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugWS, zapcore.DebugLevel),
		)
	}

	// Combine cores
	core := zapcore.NewTee(infoCore, errorCore, debugCore)

	// Create logger
	logger := zap.New(core)
	return logger
}
