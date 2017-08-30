package action

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Config is the configuration for the Action
type Config struct {
	Id       string          `json:"id"`
	Ref      string          `json:"ref"`
	Data     json.RawMessage `json:"data"`
	Metadata *ConfigMetadata `json:"metadata"`
}

// Metadata is the configuration metadata for the Action
type ConfigMetadata struct {
	Inputs  map[string]*data.Attribute `json:"inputs"`
	Outputs map[string]*data.Attribute `json:"outputs"`
}


func GetConfigInputMetadata(act Action) map[string]*data.Attribute {

	if act.Config() != nil {
		if act.Config().Metadata != nil {
			return act.Config().Metadata.Inputs
		}
	}

	return nil
}