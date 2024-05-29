package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type ErrorHook struct {
	log *logrus.Logger
}

func NewErrorHook(log *logrus.Logger) logrus.Hook {
	return &ErrorHook{log: log}
}

func (h *ErrorHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.FatalLevel, logrus.PanicLevel, logrus.ErrorLevel}
}

func (h *ErrorHook) Fire(entry *logrus.Entry) error {
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	logFileName := ""
	switch entry.Level {
	case logrus.ErrorLevel:
		logFileName = "error.log"
	case logrus.FatalLevel:
		logFileName = "fatal.log"
	case logrus.PanicLevel:
		logFileName = "panic.log"
	default:
		logFileName = "error.log"
	}

	outputPath := filepath.Join(workdir, "log", logFileName)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open or create the log file
	output, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer output.Close()

	// Set log output to the opened file
	h.log.SetOutput(output)
	defer h.log.SetOutput(h.log.Out)

	// Log the entry
	h.log.Printf("time='%v' level='%v' message='%s'\n", entry.Time, entry.Level, entry.Message)

	return nil
}
