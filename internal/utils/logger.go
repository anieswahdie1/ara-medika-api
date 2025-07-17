package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(environment string) {
	Logger = logrus.New()

	switch environment {
	case "production":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "severity",
				logrus.FieldKeyMsg:   "message",
			},
		})
	default:
		Logger.SetFormatter(&HumanFormatter{})
		Logger.SetOutput(os.Stdout)
	}

	// Set level berdasarkan env
	if environment == "development" {
		Logger.SetLevel(logrus.DebugLevel)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}
}

type HumanFormatter struct{}

func (f *HumanFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Warna untuk terminal
	color := 36 // cyan
	switch entry.Level {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = 31 // merah
	case logrus.WarnLevel:
		color = 33 // kuning
	case logrus.InfoLevel:
		color = 32 // hijau
	case logrus.DebugLevel:
		color = 36 // cyan
	}

	// Format timestamp
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")

	// Dapatkan info file dan line
	file, line := getCaller()

	// Format pesan utama
	msg := fmt.Sprintf("\x1b[%dm%-5s\x1b[0m [%s] %s:%d - %s",
		color,
		strings.ToUpper(entry.Level.String()),
		timestamp,
		file,
		line,
		entry.Message)

	// Tambahkan fields jika ada
	if len(entry.Data) > 0 {
		fields := make([]string, 0, len(entry.Data))
		for k, v := range entry.Data {
			fields = append(fields, fmt.Sprintf("\x1b[37m%s\x1b[0m=%v", k, v))
		}
		msg += " | " + strings.Join(fields, ", ")
	}

	return []byte(msg + "\n"), nil
}

func getCaller() (string, int) {
	// Skip 5 level stack:
	// 0. runtime.Caller
	// 1. getCaller
	// 2. Format
	// 3. logrus entry
	// 4. logger function call
	// 5. actual caller
	_, file, line, ok := runtime.Caller(5)
	if !ok {
		return "???", 0
	}

	// Pendekkan path file
	parts := strings.Split(file, "/")
	if len(parts) > 3 {
		file = strings.Join(parts[len(parts)-3:], "/")
	}

	return file, line
}

func SetupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}
