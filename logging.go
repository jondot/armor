package armor

import (
	"github.com/Sirupsen/logrus"
	logrus_logstash "github.com/Sirupsen/logrus/formatters/logstash"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"log"
	"log/syslog"
)

type Log struct {
	*logrus.Logger
	config *Config
}

func newLog(c *Config) *Log {
	lg := &Log{Logger: logrus.New(), config: c}
	for k, fn := range LogMap {
		if c.Exists(k) {
			fn(lg)
		}
	}

	return lg
}

type LogBuilder func(*Log)

var LogMap = map[string]LogBuilder{
	"log.console": func(lg *Log) {
		lvl, err := logrus.ParseLevel(lg.config.GetStringWithDefault("log.console.level", "debug"))
		if err != nil {
			log.Fatalf("Cannot understand log level: %s", err)
		}
		lg.Level = lvl

		var formatter logrus.Formatter = new(logrus.JSONFormatter)
		switch lg.config.GetStringWithDefault("log.console.format", "json") {
		case "text":
			formatter = new(logrus.TextFormatter)
		case "logstash":
			formatter = &logrus_logstash.LogstashFormatter{Type: lg.config.Product}
		}
		lg.Formatter = formatter
	},
	"log.syslog": func(lg *Log) {
		cfg := lg.config
		t := cfg.GetString("log.syslog.transport")
		s := cfg.GetString("log.syslog.server")
		l := lg.parseSyslog(cfg.GetString("log.hooks.syslog.level"))
		tag := cfg.Product
		hook, err := logrus_syslog.NewSyslogHook(t, s, l, tag)
		if err != nil {
			log.Fatalf("Log: cannot add syslog hook: %s", err)
		}
		lg.Hooks.Add(hook)
	},
}

func (lg *Log) parseSyslog(level string) syslog.Priority {
	switch level {
	case "panic":
		return syslog.LOG_EMERG
	case "fatal":
		return syslog.LOG_CRIT
	case "error":
		return syslog.LOG_ERR
	case "warn", "warning":
		return syslog.LOG_NOTICE
	case "info":
		return syslog.LOG_INFO
	case "debug":
		return syslog.LOG_DEBUG
	}
	return syslog.LOG_DEBUG
}
