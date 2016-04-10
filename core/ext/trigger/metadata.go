package trigger

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Metadata is the metadata for a Trigger
type Metadata struct {
	ID                string            `json:"name"`
	Settings          []*data.Attribute `json:"settings"`
	Outputs           []*data.Attribute `json:"outputs"`
	SupportsEndpoints bool              `json:"supportsEndpoints"`
	Endpoint          EndpointMetadata  `json:"endpoint"`
}

// EndpointMetadata is the metadata for a Trigger Endpoint
type EndpointMetadata struct {
	Settings []*data.Attribute `json:"settings"`
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
