package activity

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Metadata is the metadata for the Activity
type Metadata struct {
	ID      string            `json:"name"`
	Inputs  []*data.Attribute `json:"inputs"`
	Outputs []*data.Attribute `json:"outputs"`
}

// NewMetadata creates the metadata object from its json representation
func NewMetadata(jsonMetadata string) *Metadata {
	md := &Metadata{}
	err := json.Unmarshal([]byte(jsonMetadata), md)
	if err != nil {
		panic("Unable to parse activity metadata: " + err.Error())
	}

	return md
}
