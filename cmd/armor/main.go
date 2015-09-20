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
						fmt.Printf("Cannot change into directory: %s\n", app)
						os.Exit(1)
					}
					fmt.Printf("Created %s\n", app)

					// setting up project. note: godeps will be fetched on
					// user's first 'make', so no need to vendor for her here.
					os.MkdirAll("config", 0755)
					createWithTemplate(TEMPL_CONF, "config/development.yaml", b)
					createWithTemplate(TEMPL_SERVICE, "main.go", b)
					createWithTemplate(TEMPL_PROCFILE, "Procfile", b)
					createWithTemplate(TEMPL_MAKE, "Makefile", b)
					createWithTemplate(TEMPL_DOCKER, "Dockerfile", b)
					createWithTemplate(TEMPL_GITIGNORE, ".gitignore", b)

					// git related
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
