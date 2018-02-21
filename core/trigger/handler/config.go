package handler

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Config is the configuration for the Trigger Handler
type Config struct {

	parent   *Config
	Settings map[string]interface{} `json:"settings"`
	Output   map[string]interface{} `json:"output"`
	Action   *action.Config

	//for backwards compatibility
	ActionId           string           `json:"actionId"`
	ActionMappings     *data.IOMappings `json:"actionMappings,omitempty"`
	actionInputMapper  data.Mapper
	actionOutputMapper data.Mapper

	Outputs              map[string]interface{} `json:"outputs"`
	ActionOutputMappings []*data.MappingDef     `json:"actionOutputMappings,omitempty"`
	ActionInputMappings  []*data.MappingDef     `json:"actionInputMappings,omitempty"`
}
