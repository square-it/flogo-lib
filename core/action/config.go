package action

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Config is the configuration for the Action
type Config struct {
	Ref      string           `json:"ref"`
	Data     json.RawMessage  `json:"data"`
	Mappings *data.IOMappings `json:"mappings"`

	//Deprecated
	Id string `json:"id"`
	//Deprecated
	Metadata *data.IOMetadata `json:"metadata"`
	//Deprecated
	//OldData json.RawMessage
}

//func GetConfigInputMetadata(act Action) []*data.Attribute {
//
//	if act.Config() != nil {
//		if act.Config().Metadata != nil {
//			return act.Config().Metadata.Input
//		}
//	}
//
//	return nil
//}
//
//func GetConfigOutputMetadata(act Action) []*data.Attribute {
//
//	if act.Config() != nil {
//		if act.Config().Metadata != nil {
//			return act.Config().Metadata.Output
//		}
//	}
//
//	return nil
//}
