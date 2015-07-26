package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jondot/armor"
	"io/ioutil"
	"time"
)

func main() {
	arm := armor.New("producer", "1.0")
	r := arm.GinRouter()

	amqp := NewAmqp(arm.Config.GetString("amqp"), arm.Config.GetString("exchange"))
	err := amqp.Dial()
	if err != nil {
		arm.Log.Fatalf("Cannot dial amqp AMQP: %s", err)
	}

	err = amqp.Declare("usage-logs")
	if err != nil {
		arm.Log.Fatalf("Cannot declare a queue: %s", err)
	}

	go amqp.AutoHeal()

	r.GET("/", func(c *gin.Context) {
		defer arm.Metrics.Timed("amqp.produce_sample", time.Now())
		err := amqp.Publish("usage-logs", []byte("test message"))
		if err != nil {
			c.String(500, fmt.Sprintf("error: %s", err))
			return
		}

		c.String(200, "hello world")
	})

	r.POST("/", func(c *gin.Context) {
		defer arm.Metrics.Timed("amqp.produce_post_request", time.Now())
		bytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			arm.Log.Error("cannot publish", err)
		}
		amqp.Publish("usage-logs", bytes)
		c.String(200, "hello world")
	})

	arm.Run(r)
}
