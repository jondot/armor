package armor

import (
	"fmt"
	viper "github.com/jondot/viper" // my branch. see https://github.com/spf13/viper/pull/76/files
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Environment string
	Hostname    string
	Product     string
	Version     string
	Sha         string
}

func (c *Config) All() interface{} {
	return viper.AllSettings()
}

func (c *Config) Get(key string) interface{} {
	return viper.Get(key)
}

func (c *Config) GetStringArray(key string) []string {
	return viper.GetStringSlice(key)
}

func (c *Config) GetString(key string) string {
	return viper.GetString(key)
}

// viper has a bug which doesn't allow actual config to override nested default.
func (c *Config) GetStringNestedWithDefault(key string, defval string) string {
	v := viper.GetString(key)
	if v != "" {
		return v
	}
	return defval
}

func (c *Config) GetInt(key string) int {
	return viper.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c *Config) Exists(key string) bool {
	return viper.IsSet(key)
}

func newConfig(product string, ver string) *Config {
	env := pickEnv("development", "GO_ENV", "RACK_ENV", "RAILS_ENV", "NODE_ENV")
	viper.SetConfigName(env)       // name of config file (without extension)
	viper.AddConfigPath("config/") // path to look for the config file in
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(product)
	viper.AutomaticEnv()
	viper.BindEnv("port", "PORT")
	viper.SetDefault("port", "6060")
	viper.SetDefault("interface", "")
	viper.SetDefault("shafile", "REVISION")
	viper.SetDefault("middleware", []string{"request_identification", "request_tracing", "route_metrics"})

	host, _ := os.Hostname()
	viper.SetDefault("host", host)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatalf("Cannot read configuration %v", err)
	}

	host = viper.GetString("host")
	sha := readSha()

	c := &Config{
		Environment: env,
		Hostname:    host,
		Product:     product,
		Version:     ver,
		Sha:         sha,
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

func readSha() string {
	shafile := viper.GetString("shafile")
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
