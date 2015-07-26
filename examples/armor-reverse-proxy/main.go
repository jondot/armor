package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jondot/armor"
	"time"
)

func main() {
	arm := armor.New("reverseproxy", "1.0")
	r := arm.GinRouter()

	//XXX goquery
	r.GET("/", func(c *gin.Context) {
		defer arm.Metrics.Timed("timed.request", time.Now())
		arm.Metrics.Inc("foobar")
		c.String(200, "hello world")
	})

	arm.Run(r)
}
