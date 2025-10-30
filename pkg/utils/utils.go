package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	UseLogrus         bool
	LogToConsole      bool
	LogToFile         bool
	SeparateByLevel   bool
	CustomInfoFile    string
	CustomWarningFile string
	CustomErrorFile   string
	LogDirectory      string
}

type LogLevel string

const (
	InfoLevel    LogLevel = "info"
	WarningLevel LogLevel = "warning"
	ErrorLevel   LogLevel = "error"
)

// CreateLogConfigFromViper creates a LogConfig from Viper configuration
func CreateLogConfigFromViper() *LogConfig {
	loggingConfig := configuration.GetLoggingConfig()
	return &LogConfig{
		UseLogrus:         loggingConfig.UseLogrus,
		LogToConsole:      loggingConfig.LogToConsole,
		LogToFile:         loggingConfig.LogToFile,
		SeparateByLevel:   loggingConfig.LogSeparateByLevel,
		CustomInfoFile:    loggingConfig.LogCustomInfoFile,
		CustomWarningFile: loggingConfig.LogCustomWarningFile,
		CustomErrorFile:   loggingConfig.LogCustomErrorFile,
		LogDirectory:      loggingConfig.LogDirectory,
	}
}

// GetLogFileName returns the appropriate log filename based on level and custom settings
func GetLogFileName(config *LogConfig, level LogLevel) string {
	if config.SeparateByLevel {
		switch level {
		case InfoLevel:
			if config.CustomInfoFile != "" {
				return config.CustomInfoFile
			}
			return "info.log"
		case WarningLevel:
			if config.CustomWarningFile != "" {
				return config.CustomWarningFile
			}
			return "warning.log"
		case ErrorLevel:
			if config.CustomErrorFile != "" {
				return config.CustomErrorFile
			}
			return "error.log"
		default:
			return "app.log"
		}
	}
	return "app.log"
}

// CreateLogDirectory ensures the log directory exists
func CreateLogDirectory(directory string) error {
	return os.MkdirAll(directory, os.ModePerm)
}

// GetLogFileWithLevel returns a file for logging with the specified level
func GetLogFileWithLevel(config *LogConfig, level LogLevel) (*os.File, error) {
	logDir := config.LogDirectory
	if logDir == "" {
		logDir = "storage/logs"
	}

	if err := CreateLogDirectory(logDir); err != nil {
		return nil, fmt.Errorf("error creating log directory: %w", err)
	}

	filename := GetLogFileName(config, level)
	filePath := fmt.Sprintf("%s/%s", logDir, filename)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file %s: %w", filePath, err)
	}

	return file, nil
}

// GetMultiWriter returns a MultiWriter based on configuration
func GetMultiWriter(config *LogConfig, level LogLevel) (io.Writer, error) {
	var writers []io.Writer

	if config.LogToConsole {
		writers = append(writers, os.Stdout)
	}

	if config.LogToFile {
		file, err := GetLogFileWithLevel(config, level)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	return io.MultiWriter(writers...), nil
}

// LogFiber creates a Fiber logger configuration with enhanced features
// Usage: LogFiber() for default, LogFiber("custom message") for custom format
func LogFiber(message ...string) logger.Config {
	config := CreateLogConfigFromViper()

	logDir := config.LogDirectory
	if logDir == "" {
		logDir = "storage/logs"
	}

	if err := CreateLogDirectory(logDir); err != nil {
		log.Fatalf("error creating log directory: %v", err)
	}

	writer, err := GetMultiWriter(config, InfoLevel)
	if err != nil {
		log.Fatalf("error creating log writer: %v", err)
	}

	format := "[${time}]: ${ip} - ${method} ${status} ${path} ${latency} ${bytesSent}B\n"

	if len(message) > 0 && message[0] != "" {
		format = message[0]
	}

	return logger.Config{
		Output:     writer,
		Format:     format,
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}
}

// LogTerminal creates a terminal-only logger configuration
func LogTerminal() logger.Config {
	return logger.Config{
		Output:     os.Stdout,
		Format:     "[${time}]: ${ip} - ${method} ${status} ${path} ${latency} ${bytesSent}B\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}
}

// LogDev creates a development logger configuration with enhanced features
// Usage: LogDev() for default, LogDev("custom message") for custom format
func LogDev(message ...string) logger.Config {
	config := CreateLogConfigFromViper()

	writer, err := GetMultiWriter(config, ErrorLevel)
	if err != nil {
		log.Fatalf("error creating log writer: %v", err)
	}

	format := "[${time}]: ${ip} - ${method} ${status} ${path} ${latency} ${bytesSent}B\n"

	if len(message) > 0 && message[0] != "" {
		format = message[0]
	}

	return logger.Config{
		Output:     writer,
		Format:     format,
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}
}

// GetLogFile returns a file for logging with the specified filename (legacy compatibility)
func GetLogFile(config *LogConfig, filename string) *os.File {
	logDir := config.LogDirectory
	if logDir == "" {
		logDir = "storage/logs"
	}

	if err := CreateLogDirectory(logDir); err != nil {
		if config != nil && config.UseLogrus {
			logrus.Fatalf("error creating log directory: %v", err)
		} else {
			panic("error creating log directory: " + err.Error())
		}
	}

	filePath := fmt.Sprintf("%s/%s", logDir, filename)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if config != nil && config.UseLogrus {
			logrus.Fatalf("error opening log file %s: %v", filename, err)
		} else {
			panic("error opening log file " + filename + ": " + err.Error())
		}
	}

	return file
}

// ValidatorErrors converts validator errors to JSONError slice
func ValidatorErrors(err error) []response.JSONError {
	var errors []response.JSONError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorMsg := fmt.Sprintf("Field '%s' failed validation '%s'", e.Field(), e.Tag())
			if e.Param() != "" {
				errorMsg += fmt.Sprintf(" with parameter '%s'", e.Param())
			}

			errors = append(errors, response.JSONError{
				Status:  fiber.StatusBadRequest,
				Message: errorMsg,
			})
		}
	} else {
		errors = append(errors, response.JSONError{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return errors
}
