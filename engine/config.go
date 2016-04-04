package engine

import (
	"os"
	"encoding/json"
)

type Configuration struct {
	LogLevel        string                    `json:"loglevel"`
	StateServiceURI string                    `json:"state_service"`
	EngineConfig    *EngineConfig             `json:"engine"`
	Triggers        map[string]*TriggerConfig `json:"triggers"`
}

type configRep struct {
	LogLevel        string           `json:"loglevel"`
	StateServiceURI string           `json:"state_service"`
	EngineConfig    *EngineConfig    `json:"engine"`
	Triggers        []*TriggerConfig `json:"triggers"`
}

func NewConfiguration() *Configuration {
	engineConfig := NewEngineConfig()

	return &Configuration{
		LogLevel: "INFO",
		StateServiceURI:"http://localhost:9190",
		EngineConfig: engineConfig,
		Triggers:make(map[string]*TriggerConfig),
	}
}

func newConfiguration(rep *configRep) *Configuration {

	config := &Configuration{
		LogLevel:rep.LogLevel,
		StateServiceURI:rep.StateServiceURI,
		EngineConfig:rep.EngineConfig,
		Triggers:make(map[string]*TriggerConfig),
	}

	for _, trigger := range rep.Triggers {
		config.Triggers[trigger.Name] = trigger
	}

	return config
}

func LoadConfigurationFromFile(fileName string) *Configuration {

	if len(fileName) == 0 {
		panic("file name cannot be empty")
	}

	configFile, _ := os.Open(fileName)

	if configFile != nil {

		rep := &configRep{}

		decoder := json.NewDecoder(configFile)
		decodeErr := decoder.Decode(rep)
		if decodeErr != nil {
			log.Error("error:", decodeErr)
		}

		return newConfiguration(rep)
	}

	return nil
}

type TriggerConfig struct {
	Name   string `json:"name"`
	Config map[string]string  `json:"config"`
}

// EngineConfig is a configuration object used when creating an
// Engine that contains all necessary settings
type EngineConfig struct {
	NumWorkers    int `json:"workers_count"`
	WorkQueueSize int `json:"workqueue_size"`
	MaxStepCount  int `json:"stepcount_max"`
}

// NewEngineConfig creates an EngineConfig with default values
func NewEngineConfig() *EngineConfig {

	return &EngineConfig{NumWorkers: 5, WorkQueueSize: 50, MaxStepCount: 100}
}