package trigger

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Metadata is the metadata for a Trigger
type Metadata struct {
	ID     string            `json:"name"`
	Config []*data.Attribute `json:"config"`
}

// NewMetadata creates a Metadata object from the json representation
func NewMetadata(jsonMetadata string) *Metadata {
	md := &Metadata{}
	err := json.Unmarshal([]byte(jsonMetadata), md)
	if err != nil {
		panic("Unable to parse trigger metadata: " + err.Error())
	}

	return md
}
