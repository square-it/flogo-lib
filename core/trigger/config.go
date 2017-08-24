package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
)

// Config is the configuration for a Trigger
type Config struct {
	Name     string                 `json:"name"`
	Id       string                 `json:"id"`
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings"`
	Outputs  map[string]interface{} `json:"outputs"`
	Handlers []*HandlerConfig       `json:"handlers"`
}

func (c *Config) FixUp(metadata *Metadata) {

	// fix up top-level outputs
	for name, value := range c.Outputs {

		attr, ok := metadata.Outputs[name]

		if ok {
			newValue, err := data.CoerceToValue(value, attr.Type)

			if err != nil {
				//todo handle error
			} else {
				c.Outputs[name] = newValue
			}
		}
	}

	// fix up handler outputs
	for _, hc := range c.Handlers {

		hc.parent = c

		// fix up outputs
		for name, value := range hc.Outputs {

			attr, ok := metadata.Outputs[name]

			if ok {
				newValue, err := data.CoerceToValue(value, attr.Type)

				if err != nil {
					//todo handle error
				} else {
					hc.Outputs[name] = newValue
				}
			}
		}

		// create mappers
		if hc.OutputMappings != nil {
			hc.outputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: hc.OutputMappings})
		}
	}
}

func (c *Config) GetSetting(setting string) string {
	return c.Settings[setting].(string)
}

// HandlerConfig is the configuration for the Trigger Handler
type HandlerConfig struct {
	parent   *Config
	ActionId string                 `json:"actionId"`
	Settings map[string]interface{} `json:"settings"`
	Outputs  map[string]interface{} `json:"outputs"`

	OutputMappings  []*data.MappingDef `json:"outputMappings,omitempty"`
	outputMapper data.Mapper
}

func (hc *HandlerConfig) GetSetting(setting string) string {
	return hc.Settings[setting].(string)
}

func (hc *HandlerConfig) GetOutput(name string) interface{} {

	value, ok := hc.Outputs[name]

	if !ok {
		value, ok = hc.parent.Outputs[name]
	}

	return value
}

func (hc *HandlerConfig) GetOutputMapper() data.Mapper {
	return hc.outputMapper
}