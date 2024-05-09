package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

// Flags are the command line Flags
type Flags struct {
	Config string
	Debug  bool
}

/*
directories:
  shortname: foo
    path: /var/log/foo
  shortname: bar
    path: /opt/app/bar
*/

type Directory struct {
	Path string `yaml:"path"`
}

// Config contains the njmon_exporter configuration data
type Config struct {
	Interval    int                  `yaml:"scrape_seconds"`
	Directories map[string]Directory `yaml:"directories"`
	Exporter    struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	}
	Logging struct {
		Journal  bool   `yaml:"journal"`
		LevelStr string `yaml:"level"`
	} `yaml:"logging"`
}

// ParseConfig imports a yaml formatted config file into a Config struct
func ParseConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	// Define some defaults
	if config.Interval == 0 {
		config.Interval = 60
	}
	if config.Exporter.Address == "" {
		config.Exporter.Address = "0.0.0.0"
	}
	if config.Exporter.Port == 0 {
		config.Exporter.Port = 9239
	}
	return config, nil
}

// parseFlags processes arguments passed on the command line in the format
// standard format: --foo=bar
func ParseFlags() *Flags {
	f := new(Flags)
	flag.StringVar(&f.Config, "config", "dir_exporter.yml", "Path to dir_exporter configuration file")
	flag.BoolVar(&f.Debug, "debug", false, "Expand logging with Debug level messaging and format")
	flag.Parse()
	return f
}
