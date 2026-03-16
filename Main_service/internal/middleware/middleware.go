package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/baigel/lms/main-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// responseBodyWriter is a custom writer to capture the response status
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LoggerMiddleware logs POST, PUT, DELETE requests using Logrus
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read body for logging if needed
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
			latency := time.Since(start)
			statusCode := c.Writer.Status()
			
			logEntry := logger.Log.WithFields(logrus.Fields{
				"status":  statusCode,
				"latency": latency,
				"method":  method,
				"path":    c.Request.URL.Path,
				"ip":      c.ClientIP(),
			})

			if len(c.Errors) > 0 {
				logEntry.Error(c.Errors.String())
			} else if statusCode >= 400 {
				logEntry.Warn("Request failed")
			} else {
				logEntry.Info("Request handled")
			}
		}
	}
}

// ErrorHandlerMiddleware provides unified error handling
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Custom error logic can be placed here based on err.Type or err.Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message": err.Error(),
				},
			})
			return
		}
	}
}
