package app

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

// App is the configuration for the App
type Config struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Triggers    []*trigger.Config `json:"triggers"`
	Actions     []*action.Config  `json:"actions"`
}