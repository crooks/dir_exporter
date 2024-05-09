package main

import (
	"fmt"
	"gitlab/dir_exporter/config"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/crooks/jlog"
	loglevel "github.com/crooks/log-go-level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cfg   *config.Config
	flags *config.Flags
	prom  *prometheusMetrics
)

// dirSize takes a directory path and its associated shotname.  It walks the
// path and records the cumulative file sizes.
func dirSize(dirName, dirPath string) {
	startTime := time.Now()
	// Test if the requested dirPath exists
	if stat, err := os.Stat(dirPath); err == nil && stat.IsDir() {
		prom.exists.WithLabelValues(dirName).Set(1)
	} else {
		// The requested directory does not exist
		log.Debugf("%s: Directory does not exist")
		prom.exists.WithLabelValues(dirName).Set(0)
		return
	}
	var totalSize int64
	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return err
	})
	if err != nil {
		log.Warnf("Unable to parse %s: %v", dirPath, err)
		return
	}
	elapsed := time.Since(startTime).Seconds()
	prom.sizeBytes.WithLabelValues(dirName).Set(float64(totalSize))
	prom.scrapeSecs.WithLabelValues(dirName).Set(elapsed)
	log.Debugf("%s: Directory successfully parsed in %f seconds", dirPath, elapsed)
}

// dirsInfo iterates through the directories defined in the configuration and
// gathers metrics associated with each.
func dirsInfo(dirs map[string]config.Directory) {
	for k, d := range dirs {
		go dirSize(k, d.Path)
	}
}

// metricsCollector is an endless loop that periodically populates Prometheus metrics.
func metricsCollector() {
	interval := time.Duration(cfg.Interval) * time.Second
	for {
		dirsInfo(cfg.Directories)
		time.Sleep(interval)
	}
}

func main() {
	var err error
	flags = config.ParseFlags()
	cfg, err = config.ParseConfig(flags.Config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	// Define logging level and method
	loglev, err := loglevel.ParseLevel(cfg.Logging.LevelStr)
	if err != nil {
		log.Fatalf("unable to set log level: %v", err)
	}
	if cfg.Logging.Journal && jlog.Enabled() {
		log.Current = jlog.NewJournal(loglev)
	} else {
		log.Current = log.StdLogger{Level: loglev}
	}

	prom = initCollectors()
	go metricsCollector()
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%d", cfg.Exporter.Address, cfg.Exporter.Port)
	log.Infof("Dir Exporter is listening for connections on %s", exporter)
	err = http.ListenAndServe(exporter, nil)
	if err != nil {
		log.Fatalf("HTTP listener failed: %v", err)
	}
}
