package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/baigel/lms/main-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LoggerMiddleware логирует все входящие запросы.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logEntry := logger.Log.WithFields(logrus.Fields{
			"status":  statusCode,
			"latency": latency,
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"ip":      c.ClientIP(),
		})

		if len(c.Errors) > 0 {
			logEntry.Error(c.Errors.String())
		} else if statusCode >= 400 {
			logEntry.Warn("request failed")
		} else {
			logEntry.Info("request handled")
		}
	}
}

// ErrorHandlerMiddleware — единственное место, где ошибки превращаются в HTTP-ответы.
// Работает только если хендлер ещё не написал ответ (c.Writer.Written()).
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// Если хендлер уже отправил ответ — не пишем второй раз.
		// Это защита от дублирования, но в нормальном потоке сюда попадать не должны.
		if c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}

		// Неизвестная ошибка — 500, оригинал скрываем от клиента.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
