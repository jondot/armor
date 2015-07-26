package armor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jondot/armor/middleware"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)

type Armor struct {
	Config  *Config
	Log     *Log
	Metrics *Metrics
}

func New(product string, ver string) *Armor {
	m := &Armor{
		Config: newConfig(product, ver),
	}

	m.Log = newLog(m.Config)
	m.Metrics = newMetrics(m.Config)
	return m
}

func (m *Armor) GetMiddleware(name string) gin.HandlerFunc {
	switch name {
	case "request_identification":
		return middleware.GinRequestIdentification(m.Config.Hostname)
	case "request_tracing":
		return middleware.GinRequestTracing(m.Log.Logger, time.RFC3339)
	case "route_metrics":
		return middleware.GinRouteMetrics(m.Metrics.sink, m.Config.Environment, m.Config.Product, m.Config.Hostname)
	}
	return nil
}

func (m *Armor) GinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	m.useConfiguredMiddleware(r)
	ping := middleware.Ping(
		m.Config.Product,
		m.Config.Version,
		m.Config.Sha,
		m.Config.Hostname,
	)
	r.GET("/ping", func(c *gin.Context) {
		ping.ServeHTTP(c.Writer, c.Request)
	})

	return r
}

func (m *Armor) Run(g *gin.Engine) {
	m.RunWithRPC(g, nil)
}

func (m *Armor) RunWithRPC(g *gin.Engine, rpcInit func(*grpc.Server)) {
	runOn := fmt.Sprintf("%s:%s", m.Config.GetString("interface"), m.Config.GetString("port"))
	banner := `

      A R M O R
    .---.___.---.
    |     |     |   Running .....: PRODUCT_NAME
    |_____|_____|   Host/Iface ..: HOST/INTERFACE
      |___|___|
      |___|___|
     `
	// XXX replace with tmpl
	banner = strings.Replace(
		strings.Replace(
			strings.Replace(banner, "HOST", m.Config.Hostname, -1), "PRODUCT_NAME", m.Config.Product, -1),
		"INTERFACE", runOn, -1)

	log.Print(banner)
	log.Printf("-> Environment: %v", m.Config.Environment)
	log.Printf("-> Product: %v", m.Config.Product)
	log.Printf("-> Host/Interface: %v/%v", m.Config.Hostname, runOn)
	log.Printf("-> Config: %v", m.Config.All())

	if rpcInit != nil {
		rpc := newRPC(m.Config)
		rpc.Run(rpcInit)
	}
	g.Run(runOn)
}

func (m *Armor) useConfiguredMiddleware(r *gin.Engine) {
	middlewareStack := m.Config.GetStringArray("middleware")

	for _, name := range middlewareStack {
		if mid := m.GetMiddleware(name); mid != nil {
			r.Use(mid)
		}
	}
}
