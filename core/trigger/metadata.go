package trigger

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Metadata is the metadata for a Trigger
type Metadata struct {
	ID       string
	Handler  *HandlerMetadata
	Settings map[string]*data.Attribute
	Output   map[string]*data.Attribute
	Reply    map[string]*data.Attribute
}

// EndpointMetadata is the metadata for a Trigger Endpoint
type HandlerMetadata struct {
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

// UnmarshalJSON overrides the default UnmarshalJSON for Metadata
func (md *Metadata) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Name     string            `json:"name"`
		Ref      string            `json:"ref"`
		Handler  *HandlerMetadata  `json:"handler"`
		Settings []*data.Attribute `json:"settings"`
		Output   []*data.Attribute `json:"output"`
		Reply    []*data.Attribute `json:"reply"`
		//for backwards compatibility
		Endpoint *HandlerMetadata  `json:"endpoint"`
		Outputs  []*data.Attribute `json:"outputs"`
	}{}

	if err := json.Unmarshal(b, ser); err != nil {
		return err
	}

	if len(ser.Ref) > 0 {
		md.ID = ser.Ref
	} else {
		// Added for backwards compatibility
		// TODO remove and add a proper error once the BC is removed
		md.ID = ser.Name
	}

	if ser.Handler != nil {
		md.Handler = ser.Handler
	} else {
		// Added for backwards compatibility
		// TODO remove and add a proper error once the BC is removed

		if ser.Endpoint != nil {
			md.Handler = ser.Endpoint
		}
	}

	md.Settings = make(map[string]*data.Attribute, len(ser.Settings))
	md.Output = make(map[string]*data.Attribute, len(ser.Outputs))

	for _, attr := range ser.Settings {
		md.Settings[attr.Name] = attr
	}

	if len(ser.Output) > 0 {
		for _, attr := range ser.Output {
			md.Output[attr.Name] = attr
		}
	} else {
		// for backwards compatibility
		for _, attr := range ser.Outputs {
			md.Output[attr.Name] = attr
		}
	}

	return nil
}

// OutputsToAttrs converts the supplied output data to attributes
func (md *Metadata) OutputsToAttrs(outputData map[string]interface{}, coerce bool) ([]*data.Attribute, error) {

	attrs := make([]*data.Attribute, 0, len(md.Output))

	for k, a := range md.Output {
		v, _ := outputData[k]

		if coerce {
			var err error
			v, err = data.CoerceToValue(v, a.Type)

			if err != nil {
				return nil, err
			}
		}

		attrs = append(attrs, data.NewAttribute(a.Name, a.Type, v))
	}

	return attrs, nil
}
