package armor

import (
	"github.com/armon/go-metrics"
	"log"
	"strings"
	"time"
)

type Metrics struct {
	sink   metrics.MetricSink
	config *Config
	Prefix string
}

func newMetrics(c *Config) *Metrics {
	fan := metrics.FanoutSink{}
	for k, fn := range MetricsMap {
		if c.Exists(k) {
			fan = append(fan, fn(c))
		}
	}

	pref := c.GetStringNestedWithDefault("metrics.prefix", c.Environment)
	metrics.NewGlobal(metrics.DefaultConfig(strings.Join([]string{pref, c.Product, "health"}, ".")), fan)
	return &Metrics{
		sink:   fan,
		config: c,
		Prefix: pref,
	}
}

func (m *Metrics) Gauge(k string, f float32) {
	m.sink.SetGauge(m.makeMetricKey(k), f)
}

func (m *Metrics) Inc(k string) {
	m.sink.IncrCounter(m.makeMetricKey(k), 1)
}

func (m *Metrics) Time(k string, microSecs float32) {
	m.sink.AddSample(m.makeMetricKey(k), microSecs)
}

func (m *Metrics) Timed(k string, start time.Time) {
	m.sink.AddSample(m.makeMetricKey(k), float32(time.Now().Sub(start)/time.Microsecond))
}

func (m *Metrics) makeMetricKey(k string) []string {
	c := m.config
	return []string{m.Prefix, c.Product, k, c.Hostname}
}

type MetricsBuilder func(*Config) metrics.MetricSink

var MetricsMap = map[string]MetricsBuilder{
	"metrics.inmem": func(c *Config) metrics.MetricSink {
		interval := c.GetInt("metrics.inmem.interval_secs")
		retain := c.GetInt("metrics.inmem.retain_mins")
		inm := metrics.NewInmemSink(time.Duration(interval)*time.Second, time.Duration(retain)*time.Minute)
		metrics.DefaultInmemSignal(inm)
		return inm
	},
	"metrics.statsd": func(c *Config) metrics.MetricSink {
		statsd, err := metrics.NewStatsdSink(c.GetString("metrics.statsd.server"))
		if err != nil {
			log.Fatalf("Cannot hookup statsd")
		}
		return statsd
	},
}
