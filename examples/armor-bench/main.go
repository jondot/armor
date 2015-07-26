package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jondot/armor"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(4)
	arm := armor.New("bench", "1.0")
	r := arm.GinRouter()

	r.GET("/", func(c *gin.Context) {
		defer arm.Metrics.Timed("timed.request", time.Now())
		arm.Metrics.Inc("bench")
		c.String(200, "bench me")
	})

	arm.Run(r)
}
