package armor

import (
	"fmt"
	"github.com/olebedev/config"

	"io/ioutil"
	"log"
	"os"
	"path"
)

type Config struct {
	Environment string
	Hostname    string
	Product     string
	Version     string
	Sha         string
	cfg         *config.Config
}

func (c *Config) All() string {
	r, err := config.RenderYaml(c.cfg)
	if err != nil {
		log.Fatalf("Cannot render configuration: %v", err)
	}
	return r
}

func (c *Config) Map(key string) map[string]interface{} {
	return c.cfg.UMap(key)
}

func (c *Config) GetStringArray(key string) []string {
	list := c.cfg.UList(key)
	new := make([]string, len(list))
	for i, v := range list {
		s, err := fmtString(v)
		if err != nil {
			log.Fatalf("Configuration: cannot get a string array for %s. Error: %s", err)
		}
		new[i] = s
	}
	return new
}

func fmtString(n interface{}) (string, error) {
	switch n := n.(type) {
	case bool, float64, int:
		return fmt.Sprint(n), nil
	case string:
		return n, nil
	}
	return "", fmt.Errorf("expected this value to be convertible to string: %v")
}

func (c *Config) GetString(key string) string {
	return c.cfg.UString(key)
}

func (c *Config) GetStringWithDefault(key string, defval string) string {
	return c.cfg.UString(key, defval)
}

func (c *Config) GetInt(key string) int {
	return c.cfg.UInt(key)
}

func (c *Config) GetBool(key string) bool {
	return c.cfg.UBool(key)
}

func (c *Config) Exists(key string) bool {
	_, err := c.cfg.Get(key)
	return err == nil
}

func newConfig(product string, ver string) *Config {
	env := pickEnv("development", "GO_ENV", "RACK_ENV", "RAILS_ENV", "NODE_ENV")

	cfg, err := config.ParseYamlFile(path.Join("config", env+".yaml"))

	if err != nil {
		log.Fatalf("Cannot read configuration: %v", err)
	}

	cfg.Set("interface", cfg.UString("interface", ""))
	cfg.Set("port", cfg.UString("port", "6060"))
	cfg.Set("shafile", cfg.UString("shafile", "REVISION"))
	cfg.Set("middleware",
		cfg.UList("middleware",
			[]interface{}{"request_identification", "request_tracing", "route_metrics"}))

	host, _ := os.Hostname()
	cfg.Set("host", cfg.UString("host", host))

	host = cfg.UString("host")


	// now override any of the defaults / configs with ENV
	cfg = cfg.Env()


	sha := readSha(cfg.UString("shafile"))

	c := &Config{
		Environment: env,
		Hostname:    host,
		Product:     product,
		Version:     ver,
		Sha:         sha,
		cfg:         cfg,
	}
	return c
}

func pickEnv(defval string, opts ...string) string {
	for _, opt := range opts {
		v := os.Getenv(opt)
		if v != "" {
			return v
		}
	}

	return defval
}

func readSha(shafile string) string {
	bytes, err := ioutil.ReadFile(shafile)
	return conditionally(err == nil, fmt.Sprintf("%s", bytes), "no-rev").(string)
}

//
// Don't confuse yourself - this is not really an if construct.
// On call site arguments are immediately evaluated.
func conditionally(test bool, prime interface{}, alt interface{}) interface{} {
	if test {
		return prime
	} else {
		return alt
	}
}
