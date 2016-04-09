package engine

import (
	"encoding/json"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
)

// Config is the configuration for the engine
type Config struct {
	LogLevel        string                     `json:"loglevel"`
	StateServiceURI string                     `json:"stateServiceURI"`
	Triggers        map[string]*trigger.Config `json:"triggers"`
	RunnerConfig    *RunnerConfig              `json:"triggers"`
	TesterConfig    *TesterConfig              `json:"tester,omitempty"`
}

// RunnerConfig is the configuration for the engine level runner
type RunnerConfig struct {
	Type   string               `json:"type"`
	Pooled *runner.PooledConfig `json:"pooled,omitempty"`
	Direct *runner.DirectConfig `json:"direct,omitempty"`
}

// TesterConfig is the configuration for the engine tester
type TesterConfig struct {
	Enabled  bool              `json:"enabled"`
	Settings map[string]string `json:"settings"`
}

type serEngineConfig struct {
	LogLevel        string            `json:"loglevel"`
	StateServiceURI string            `json:"stateServiceURI"`
	Triggers        []*trigger.Config `json:"triggers"`
	RunnerConfig    *RunnerConfig     `json:"processRunner"`
	TesterConfig    *TesterConfig     `json:"tester,omitempty"`
}

// DefaultConfig returns the default engine configuration
func DefaultConfig() *Config {

	var engineConfig Config

	engineConfig.LogLevel = "DEBUG"
	engineConfig.Triggers = make(map[string]*trigger.Config)
	engineConfig.RunnerConfig = defaultRunnerConfig()
	engineConfig.TesterConfig = defaultTesterConfig()

	return &engineConfig
}

// MarshalJSON marshals the EngineConfig to JSON
func (ec *Config) MarshalJSON() ([]byte, error) {

	var triggers []*trigger.Config

	for _, value := range ec.Triggers {
		triggers = append(triggers, value)
	}

	return json.Marshal(&serEngineConfig{
		LogLevel:        ec.LogLevel,
		StateServiceURI: ec.StateServiceURI,
		Triggers:        triggers,
		RunnerConfig:    ec.RunnerConfig,
		TesterConfig:    ec.TesterConfig,
	})
}

// UnmarshalJSON unmarshals EngineConfog from JSON
func (ec *Config) UnmarshalJSON(data []byte) error {

	ser := &serEngineConfig{}
	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	ec.LogLevel = ser.LogLevel
	ec.StateServiceURI = ser.StateServiceURI

	if ser.RunnerConfig != nil {
		ec.RunnerConfig = ser.RunnerConfig
	} else {
		ec.RunnerConfig = defaultRunnerConfig()
	}

	if ser.TesterConfig != nil {
		ec.TesterConfig = ser.TesterConfig
	} else {
		ec.TesterConfig = defaultTesterConfig()
	}

	ec.Triggers = make(map[string]*trigger.Config)

	for _, value := range ser.Triggers {
		ec.Triggers[value.Name] = value
	}

	return nil
}

//LoadConfigFromFile loads the engine Config from the specified JSON file
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

func defaultTesterConfig() *TesterConfig {
	return &TesterConfig{Enabled: true, Settings: map[string]string{"port": "8080"}}
}

func defaultRunnerConfig() *RunnerConfig {
	return &RunnerConfig{Type: "pooled", Pooled: &runner.PooledConfig{NumWorkers: 5, WorkQueueSize: 50, MaxStepCount: 100}}
}
