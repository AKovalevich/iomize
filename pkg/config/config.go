package config

import (
	"sync"
	"time"

	"github.com/AKovalevich/iomize/pkg/entrypoint"
	"github.com/AKovalevich/iomize/pkg/pipeline"
)

const (
	// DefaultGraceTimeout controls how long Scrabbler serves pending requests
	// prior to shutting down.
	DefaultGraceTimeout = 10 * time.Second

	// DefaultConfigFileName path to configuration file
	DefaultConfigPath = "configuration.default.toml"
)

// The main Scrabbler configuration
type MainConfiguration struct {
	sync.RWMutex
	// Main configuration
	Debug          bool   `yaml:"debug"`
	LogLevel       string `yaml:"log_level"`
	ConfigFilePath string `yaml:"config_file_path"`
	// Scrabbler server configuration
	Port string `yaml:"port"`
	Host string `yaml:"host"`
	// Shutdown configuration
	GraceTimeOut   int    `yaml:"grace_time_out"`
	PipeConfigPath string `yaml:"pipe_config_path"`
	EntryPoints    entrypoint.EntrypointList
	PipeLineList   pipeline.PipeLineList
}

func NewConfiguration() *MainConfiguration {
	return &MainConfiguration{}
}

// Reload configuration
func (config *MainConfiguration) Reload() {}
