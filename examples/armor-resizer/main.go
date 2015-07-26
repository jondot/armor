package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jondot/armor"
	"gopkg.in/h2non/bimg.v0"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	arm := armor.New("resizer", "1.0")
	r := arm.GinRouter()

	client := &http.Client{}

	r.GET("/", func(c *gin.Context) {
		defer arm.Metrics.Timed("timed.request", time.Now())
		arm.Metrics.Inc("foobar")
		c.String(200, "hello world")
	})

	r.GET("/json", func(c *gin.Context) {
		start := time.Now()
		resp, err := client.Get("https://upload.wikimedia.org/wikipedia/commons/c/cb/Pocket-Gopher_Ano-Nuevo-SP.jpg")
		arm.Metrics.Timed("image.fetch", start)
		if err != nil {
			c.JSON(200, gin.H{"error": "true"})
			return
		}

		buffer, _ := ioutil.ReadAll(resp.Body)
		newImage, _ := bimg.NewImage(buffer).Resize(800, 600)
		arm.Metrics.Timed("image.resize", start)

		c.Data(200, "image/jpeg", newImage)

	})

	r.GET("/xml", func(c *gin.Context) {
		c.XML(200, gin.H{"message": "hey", "status": 200})
	})

	arm.Run(r)
}
