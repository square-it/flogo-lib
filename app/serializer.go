package app

import (
	"os"
	"encoding/json"
)

const (
	FLOGO_CONFIG_PATH = "flogo.json"
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
	flogo, err := os.Open(FLOGO_CONFIG_PATH)
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



