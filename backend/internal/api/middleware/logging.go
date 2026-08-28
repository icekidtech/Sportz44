package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logging logs structured request metadata for every request.
func Logging(l *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		l.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  c.Writer.Status(),
			"latency": time.Since(start).String(),
			"ip":      c.ClientIP(),
		}).Info("request")
	}
}
