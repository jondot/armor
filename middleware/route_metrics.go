package middleware

import (
	"github.com/armon/go-metrics"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func GinRouteMetrics(mtr metrics.MetricSink, env string, product string, host string) gin.HandlerFunc {
	const (
		routes = "routes"
		slash  = "/"
		under  = "_"
	)
	return func(c *gin.Context) {
		r := c.Request

		// Start timer
		start := time.Now()
		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		mtr.AddSample(
			[]string{
				env,
				product,
				routes,
				strings.Replace(r.URL.Path, slash, under, -1),
				r.Method,
				host,
			},
			float32(end.Sub(start)),
		)
	}
}
