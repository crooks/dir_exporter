package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix string = "dir"
)

type prometheusMetrics struct {
	exists     *prometheus.GaugeVec
	sizeBytes  *prometheus.GaugeVec
	scrapeSecs *prometheus.GaugeVec
}

func addPrefix(s string) string {
	return fmt.Sprintf("%s_%s", prefix, s)
}

func initCollectors() *prometheusMetrics {
	dir := new(prometheusMetrics)

	dir.exists = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("exists"),
			Help: "Does the requested directory exist",
		},
		[]string{"dir"},
	)
	prometheus.MustRegister(dir.exists)

	dir.sizeBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("size_bytes"),
			Help: "The size of a given directory in Bytes",
		},
		[]string{"dir"},
	)
	prometheus.MustRegister(dir.sizeBytes)

	dir.scrapeSecs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("scrape_secs"),
			Help: "The elapsed time to calculate the directory size",
		},
		[]string{"dir"},
	)
	prometheus.MustRegister(dir.scrapeSecs)

	return dir
}
