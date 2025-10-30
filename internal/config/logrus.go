package config

import (
	"fmt"
	"io"
	"os"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/sirupsen/logrus"
)

// LogrusLevelHook is a hook that writes different log levels to different files
type LogrusLevelHook struct {
	config  *utils.LogConfig
	writers map[logrus.Level]io.Writer
}

// NewLogrusLevelHook creates a new LogrusLevelHook
func NewLogrusLevelHook(config *utils.LogConfig) (*LogrusLevelHook, error) {
	hook := &LogrusLevelHook{
		config:  config,
		writers: make(map[logrus.Level]io.Writer),
	}

	if config.SeparateByLevel {
		infoWriter, err := utils.GetMultiWriter(config, utils.InfoLevel)
		if err != nil {
			return nil, fmt.Errorf("failed to create info writer: %w", err)
		}
		hook.writers[logrus.InfoLevel] = infoWriter
		hook.writers[logrus.DebugLevel] = infoWriter
		hook.writers[logrus.TraceLevel] = infoWriter

		warningWriter, err := utils.GetMultiWriter(config, utils.WarningLevel)
		if err != nil {
			return nil, fmt.Errorf("failed to create warning writer: %w", err)
		}
		hook.writers[logrus.WarnLevel] = warningWriter

		errorWriter, err := utils.GetMultiWriter(config, utils.ErrorLevel)
		if err != nil {
			return nil, fmt.Errorf("failed to create error writer: %w", err)
		}
		hook.writers[logrus.ErrorLevel] = errorWriter
		hook.writers[logrus.FatalLevel] = errorWriter
		hook.writers[logrus.PanicLevel] = errorWriter
	} else {
		allLevelWriter, err := utils.GetMultiWriter(config, utils.InfoLevel)
		if err != nil {
			return nil, fmt.Errorf("failed to create writer: %w", err)
		}

		for _, level := range logrus.AllLevels {
			hook.writers[level] = allLevelWriter
		}
	}

	return hook, nil
}

// Levels returns the levels this hook is interested in
func (hook *LogrusLevelHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire writes the log entry to the appropriate writer based on level
func (hook *LogrusLevelHook) Fire(entry *logrus.Entry) error {
	writer, exists := hook.writers[entry.Level]
	if !exists {
		writer = hook.writers[logrus.InfoLevel]
	}

	if writer != nil {
		formatter := &logrus.TextFormatter{
			FullTimestamp:    true,
			DisableColors:    false,
			DisableTimestamp: false,
			TimestampFormat:  "02-Jan-2006 15:04:05",
		}

		bytes, err := formatter.Format(entry)
		if err != nil {
			return err
		}

		_, err = writer.Write(bytes)
		return err
	}

	return nil
}

// SetupLogRus creates and configures a Logrus logger with enhanced features
func SetupLogRus(config *utils.LogConfig) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.TraceLevel)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableColors:    false,
		DisableTimestamp: false,
		TimestampFormat:  "02-Jan-2006 15:04:05",
	})

	if config.LogToConsole && !config.LogToFile {
		log.SetOutput(os.Stdout)
		return log
	}

	if config.SeparateByLevel && config.LogToFile {
		if !config.LogToConsole {
			log.SetOutput(io.Discard)
		} else {
			log.SetOutput(os.Stdout)
		}

		hook, err := NewLogrusLevelHook(config)
		if err != nil {
			logrus.Fatalf("Failed to create logrus hook: %v", err)
		}
		log.AddHook(hook)
	} else {
		writer, err := utils.GetMultiWriter(config, utils.InfoLevel)
		if err != nil {
			logrus.Fatalf("Failed to create log writer: %v", err)
		}
		log.SetOutput(writer)
	}

	return log
}
