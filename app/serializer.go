package app

import (
	"os"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/config"
)

// defaultSerializer implementation of AppSerializer
type defaultSerializer struct {
}

// AppSerializer interface to implement to provide the app configuration
type AppSerializer interface {
	GetApp() (*Config, error)
}

// DefaultSerializer returns the default App Serializer
func DefaultSerializer() AppSerializer {
	return &defaultSerializer{}
}

// GetApp returns the app configuration
func (d *defaultSerializer) GetApp() (*Config, error){

	configPath := config.GetFlogoConfigPath()

	flogo, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(flogo)
	app := &Config{}
	err = jsonParser.Decode(&app)
	if err != nil {
		return nil, err
	}

	return app, nil
}



