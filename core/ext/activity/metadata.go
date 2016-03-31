package activity

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"encoding/json"
)

type Metadata struct {
	ID string `json:"name"`
	Inputs []*data.Attribute `json:"inputs"`
	Outputs []*data.Attribute `json:"outputs"`
}

func NewMetadata(jsonMetadata string) *Metadata {
	md := &Metadata{}
	err := json.Unmarshal([]byte(jsonMetadata), md)
	if err != nil {
		panic("Unable to parse activity metadata: " + err.Error())
	}

	return md
}