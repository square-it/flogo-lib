package engine

import (
	"encoding/json"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/service"
)

// Config is the configuration for the engine
type Config struct {
	LogLevel     string                     `json:"loglevel"`
	RunnerConfig *RunnerConfig              `json:"processRunner"`
	Triggers     map[string]*trigger.Config `json:"triggers"`
	Services     map[string]*service.Config `json:"services"`
}

// RunnerConfig is the configuration for the engine level runner
type RunnerConfig struct {
	Type   string               `json:"type"`
	Pooled *runner.PooledConfig `json:"pooled,omitempty"`
	Direct *runner.DirectConfig `json:"direct,omitempty"`
}

type serEngineConfig struct {
	LogLevel     string            `json:"loglevel"`
	RunnerConfig *RunnerConfig     `json:"processRunner"`
	Triggers     []*trigger.Config `json:"triggers"`
	Services     []*service.Config `json:"services"`
}

// DefaultConfig returns the default engine configuration
func DefaultConfig() *Config {

	var engineConfig Config

	engineConfig.LogLevel = "DEBUG"
	engineConfig.Triggers = make(map[string]*trigger.Config)
	engineConfig.RunnerConfig = defaultRunnerConfig()
	engineConfig.Services = service.DefaultServicesConfig()

	return &engineConfig
}

// MarshalJSON marshals the EngineConfig to JSON
func (ec *Config) MarshalJSON() ([]byte, error) {

	var triggers []*trigger.Config

	for _, value := range ec.Triggers {
		triggers = append(triggers, value)
	}

	var services []*service.Config

	for _, value := range ec.Services {
		services = append(services, value)
	}

	return json.Marshal(&serEngineConfig{
		LogLevel:     ec.LogLevel,
		RunnerConfig: ec.RunnerConfig,
		Triggers:     triggers,
		Services:     services,
	})
}

// UnmarshalJSON unmarshals EngineConfog from JSON
func (ec *Config) UnmarshalJSON(data []byte) error {

	ser := &serEngineConfig{}
	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	ec.LogLevel = ser.LogLevel

	if ser.RunnerConfig != nil {
		ec.RunnerConfig = ser.RunnerConfig
	} else {
		ec.RunnerConfig = defaultRunnerConfig()
	}

	if ser.Services != nil {
		ec.Services = make(map[string]*service.Config)

		for _, value := range ser.Services {
			ec.Services[value.Name] = value
		}
	} else {
		ec.Services = service.DefaultServicesConfig()
	}

	ec.Triggers = make(map[string]*trigger.Config)

	for _, value := range ser.Triggers {
		ec.Triggers[value.Name] = value
	}

	return nil
}

// LoadConfigFromFile loads the engine Config from the specified JSON file
func LoadConfigFromFile(fileName string) *Config {

	if len(fileName) == 0 {
		panic("file name cannot be empty")
	}

	configFile, _ := os.Open(fileName)

	if configFile != nil {

		engineConfig := &Config{}

		decoder := json.NewDecoder(configFile)
		decodeErr := decoder.Decode(engineConfig)
		if decodeErr != nil {
			log.Error("error:", decodeErr)
		}

		return engineConfig
	}

	return nil
}

// LoadConfigFromJSON loads the engine Config from the specified JSON file
func LoadConfigFromJSON(configJSON string) *Config {

	engineConfig := &Config{}
	decodeErr := json.Unmarshal([]byte(configJSON), engineConfig)
	if decodeErr != nil {
		log.Error("error:", decodeErr)
	}

	return engineConfig
}

func defaultRunnerConfig() *RunnerConfig {
	return &RunnerConfig{Type: "pooled", Pooled: &runner.PooledConfig{NumWorkers: 5, WorkQueueSize: 50, MaxStepCount: 100}}
}
