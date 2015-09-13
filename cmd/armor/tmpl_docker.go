package main

const TEMPL_DOCKER = `
FROM tianon/true
COPY ./config /
COPY ./snapshot/linux_amd64/* /
EXPOSE 6060
CMD ["/{{.Product}}"]
`
const TEMPL_PROCFILE = `web: {{.Product}}
`
const TEMPL_SERVICE = `
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jondot/armor"
	"time"
)

func main() {
	m := armor.New("{{.Product}}", "1.0")
	r := m.GinRouter()

	r.GET("/", func(c *gin.Context) {
		defer m.Metrics.Timed("timed.request", time.Now())
		m.Metrics.Inc("foobar")
		c.String(200, "hello world")
	})

	m.Run(r)
}
`
