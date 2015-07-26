package middleware

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"time"
)

func GinRequestTracing(log *logrus.Logger, timeFormat string) gin.HandlerFunc {
	const (
		status    = "status"
		method    = "method"
		path      = "path"
		ip        = "ip"
		latency   = "latency"
		useragent = "user-agent"
		requestid = "request-id"
		tm        = "time"
		xreq      = "X-Request-Id"
		errors    = "errors"
	)
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		// XXX short circuit if not using correct log level
		// but dont forget errors

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()

		clientIP := c.ClientIP()
		// XXX const strings
		e := log.WithFields(logrus.Fields{
			status:    w.Status(),
			method:    r.Method,
			path:      r.URL.Path,
			ip:        clientIP,
			latency:   end.Sub(start),
			useragent: r.UserAgent(),
			requestid: w.Header().Get(xreq),
			tm:        end.Format(timeFormat),
		})
		if len(c.Errors) > 0 {
			e.WithField(errors, c.Errors.String()).Error(r.URL)
		} else {
			e.Info(r.URL)
		}
	}
}
