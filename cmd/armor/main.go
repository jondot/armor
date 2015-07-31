package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"text/template"
)

type bucket struct {
	Product string
}

func main() {
	app := cli.NewApp()
	app.Name = "armor"
	app.Usage = "armor [command]"
	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "Creates an Armor service",
			Action: func(c *cli.Context) {
				if app := c.Args().First(); app != "" {

					b := bucket{Product: app}

					os.MkdirAll(app, 0755)
					if err := os.Chdir(app); err != nil {
						fmt.Printf("Cannot change into directory: %s", app)
						os.Exit(1)
					}
					fmt.Printf("Created %s\n", app)

					os.MkdirAll("config", 0755)

					createWithTemplate(TEMPL_CONF, "config/development.yaml", b)
					createWithTemplate(TEMPL_SERVICE, "main.go", b)
					createWithTemplate(TEMPL_PROCFILE, "Procfile", b)
					createWithTemplate(TEMPL_MAKE, "Makefile", b)
					createWithTemplate(TEMPL_DOCKER, "Dockerfile", b)

					fmt.Print(shellExec("git", "init"))
					fmt.Print(shellExec("git", "add", "-A"))
					fmt.Print(shellExec("git", "commit", "-am", "first commit"))
				}
			},
		},
	}

	app.Run(os.Args)
}

func createWithTemplate(tmpltext string, fname string, b bucket) {
	tmpl, err := template.New(fname).Parse(tmpltext)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fname)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	tmpl.Execute(f, b)
	fmt.Printf("Wrote %s\n", fname)
}

func shellExec(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Printf("Cannot run shell command: %s %s, error: %v (%s)", cmd, args, err, out)
		os.Exit(1)
	}
	return string(out)
}

const TEMPL_MAKE = `
default:
	@godep go build
	@ls -ltrh

setup: goxc.ok
	@echo Installing dependency management tools...
	go get github.com/tools/godep

goxc.ok:
	@echo Installing crossbuild tooling. This will take a while...
	go get github.com/laher/goxc
	goxc -t
	touch goxc.ok

heroku:
	@echo Bootstrapping with godep
	@go get github.com/tools/godep
	@godep save
	@git add -A .
	@git commit -am "dependencies"
	@echo Creating a Heroku Go app...
	@heroku create -b https://github.com/kr/heroku-buildpack-go.git
	@git push heroku master

test:
	go test

release:
	godep save
	goxc -env GOPATH=` + "`godep path`" + ` -bc="linux,amd64" -d . xc # we only use basic xc for now, see github.com/laher/goxc for more

docker: release
	@docker build -t {{.Product}} .
	@echo Container [{{.Product}}] built. Run with: make docker-run

docker-run:
	docker run -p 80:6060 {{.Product}}

.PHONY: heroku build test setup release docker docker-run
`

const TEMPL_DOCKER = `
FROM scratch
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
const TEMPL_CONF = `

#====================================================================
# Your own hierarchical config
#
# Examples:
#   arm.Config.GetInt("foo.bar")
#   arm.Config.GetString("foo.baz")
#====================================================================
foo:
  bar: 1
  baz: yep!



#====================================================================
# Service
#
# Below are service level / operational knobs.
#====================================================================
#shafile: REVISION
#port: 6060
#interface: 0.0.0.0



#====================================================================
# RPC - we do RPC via gRPC and protocolbuffers (http://grpc.io)
#
# Use RPC in addition to HTTP in cases where you want service-to-
# service communication (microservices), oplog/auditlog, or need to
# handle specialized performance-oriented use cases.
#====================================================================
#rpc:
#  port: 41555
#  interface: 0.0.0.0



#====================================================================
# Middleware
#
# These are the list of prebaked middleware, you can mix and match
# by listing or not listing them here.
#====================================================================
middleware:
  - request_identification
  - request_tracing
  - route_metrics



#====================================================================
# Metrics
# You can have metrics sent to multiple sinks.
#
# See github.com/armon/go-metrics for types of sinks
#====================================================================
#metrics:
#  prefix: dev
#  inmem:
#    interval_secs: 10
#    retain_mins: 1
#  statsd:
#    server: localhost:514



#====================================================================
# Logging
#
# Along with custom hooks, you can config formatting, levels and
# where to dump more critical errors via hooks.
#====================================================================
log:
  console:
    level: error
#  hooks:
#    syslog:
#      transport: udp
#      server: localhost:514
#      level: debug


`
