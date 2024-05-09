package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	testFileName := "dir_exporter.cfg"
	testCfgFile, err := os.CreateTemp("", testFileName)
	if err != nil {
		t.Fatalf("Unable to create config test file: %v", err)
	}
	defer os.Remove(testCfgFile.Name())
	_, err = testCfgFile.WriteString(`---
directories:
  foo:
    path: /var/tmp/foo
  bar:
    path: /opt/app/bar
logging:
  journal: true
  level: debug
`)
	if err != nil {
		t.Fatalf("Unable to write to config test file: %v", err)
	}
	testCfgFile.Close()

	cfg, err := ParseConfig(testCfgFile.Name())
	if err != nil {
		t.Fatalf("ParseConfig failed with: %v", err)
	}
	if !cfg.Logging.Journal {
		t.Fatal("cfg.logging.journal should be true")
	}
	if cfg.Directories["foo"].Path != "/var/tmp/foo" {
		t.Fatalf("Expected /var/tmp/foo but got: %s", cfg.Directories["foo"].Path)
	}
	if cfg.Directories["bar"].Path != "/opt/app/bar" {
		t.Fatalf("Expected /opt/app/bar but got: %s", cfg.Directories["foo"].Path)
	}
}

func TestFlags(t *testing.T) {
	f := ParseFlags()
	expectingConfig := "njmon_exporter.yml"
	if f.Config != expectingConfig {
		t.Fatalf("Unexpected config flag: Expected=%s, Got=%s", expectingConfig, f.Config)
	}
	if f.Debug {
		t.Fatal("Unexpected debug flag: Expected=false, Got=true")
	}
}
