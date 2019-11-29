package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	configCheckCountMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "config_check_count",
		Help: "Count of how many times telegraf checked its config.",
	})
)

func init() {
	prometheus.MustRegister(configCheckCountMetric)
}
