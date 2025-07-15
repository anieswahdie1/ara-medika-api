package utils

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Buat response writer wrapper untuk capture status code
			ww := &responseWriterWrapper{ResponseWriter: w}

			// Eksekusi handler
			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			Logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     ww.status,
				"duration":   duration.String(),
				"ip":         r.RemoteAddr,
				"referer":    r.Referer(),
				"user_agent": r.UserAgent(),
			}).Info("HTTP request")
		})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterWrapper) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Helper untuk log dengan context
func LogWithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	// Ambil request ID dari context jika ada
	if reqID, ok := ctx.Value("request_id").(string); ok {
		fields["request_id"] = reqID
	}

	// Tambahkan field lain dari context jika diperlukan
	// ...

	return Logger.WithFields(fields)
}

// Helper untuk log error dengan stack trace
func LogErrorWithTrace(err error, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(logrus.Fields)
	}

	// Tambahkan stack trace untuk environment development
	if Logger.Level == logrus.DebugLevel {
		fields["stack"] = string(debug.Stack())
	}

	Logger.WithFields(fields).WithError(err).Error(message)
}
