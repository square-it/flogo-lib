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
	Output   map[string]interface{} `json:"output"`
	Handlers []*HandlerConfig       `json:"handlers"`

	//for backwards compatibility
	Outputs map[string]interface{} `json:"outputs"`
}

func (c *Config) FixUp(metadata *Metadata) {

	//for backwards compatibility
	if len(c.Output) == 0 {
		c.Output = c.Outputs
	}

	// fix up top-level outputs
	for name, value := range c.Output {

		attr, ok := metadata.Output[name]

		if ok {
			newValue, err := data.CoerceToValue(value, attr.Type)

			if err != nil {
				//todo handle error
			} else {
				c.Output[name] = newValue
			}
		}
	}

	// fix up handler outputs
	for _, hc := range c.Handlers {

		hc.parent = c

		//for backwards compatibility
		if len(hc.Output) == 0 {
			hc.Output = hc.Outputs
		}

		// fix up outputs
		for name, value := range hc.Output {

			attr, ok := metadata.Output[name]

			if ok {
				newValue, err := data.CoerceToValue(value, attr.Type)

				if err != nil {
					//todo handle error
				} else {
					hc.Output[name] = newValue
				}
			}
		}

		// create mappers
		if hc.ActionMappings != nil {
			if hc.ActionMappings.Input != nil {
				hc.actionInputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: hc.ActionMappings.Input})
			}
			if hc.ActionMappings.Output != nil {
				hc.actionOutputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: hc.ActionMappings.Output})
			}
		} else {
			if hc.ActionInputMappings != nil {
				hc.actionInputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: hc.ActionInputMappings})
			}
			if hc.ActionOutputMappings != nil {
				hc.actionOutputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: hc.ActionOutputMappings})
			}
		}
	}
}

func (c *Config) GetSetting(setting string) string {
	return c.Settings[setting].(string)
}

// HandlerConfig is the configuration for the Trigger Handler
type HandlerConfig struct {
	parent             *Config
	ActionId           string                 `json:"actionId"`
	Settings           map[string]interface{} `json:"settings"`
	Output             map[string]interface{} `json:"output"`
	ActionMappings     *Mappings `json:"actionMappings,omitempty"`
	actionInputMapper  data.Mapper
	actionOutputMapper data.Mapper

	//for backwards compatibility
	Outputs              map[string]interface{} `json:"outputs"`
	ActionOutputMappings []*data.MappingDef `json:"actionOutputMappings,omitempty"`
	ActionInputMappings  []*data.MappingDef `json:"actionInputMappings,omitempty"`
}

type Mappings struct {
	Input  []*data.MappingDef `json:"input,omitempty"`
	Output []*data.MappingDef `json:"output,omitempty"`
}

func (hc *HandlerConfig) GetSetting(setting string) string {
	return hc.Settings[setting].(string)
}

func (hc *HandlerConfig) GetOutput(name string) (interface{}, bool) {

	value, exists := hc.Output[name]

	if !exists {
		value, exists = hc.parent.Output[name]
	}

	return value, exists
}

func (hc *HandlerConfig) GetActionOutputMapper() data.Mapper {
	return hc.actionOutputMapper
}

func (hc *HandlerConfig) GetActionInputMapper() data.Mapper {
	return hc.actionInputMapper
}

type TriggerActionInputGenerator struct {
	handlerConfig  *HandlerConfig
	triggerOutputs []*data.Attribute
}

func NewTriggerActionInputGenerator(metadata *Metadata, config *HandlerConfig, outputs map[string]interface{}) *TriggerActionInputGenerator {

	outAttrs := metadata.Output

	attrs := make([]*data.Attribute, 0, len(outAttrs))

	for name, outAttr := range outAttrs {
		value, exists := outputs[name]

		if !exists {
			value, exists = config.GetOutput(name)
		}

		//todo if complex_object, handle referenced metadata

		if exists {
			attrs = append(attrs, data.NewAttribute(name, outAttr.Type, value))
		}
	}

	return &TriggerActionInputGenerator{handlerConfig: config, triggerOutputs: attrs}
}

func (ig *TriggerActionInputGenerator) GenerateInputs(inputMetadata []*data.Attribute) map[string]interface{} {

	outputMapper := ig.handlerConfig.GetActionOutputMapper()
	inScope := data.NewSimpleScope(ig.triggerOutputs, nil)
	outScope := data.NewFixedScope(inputMetadata)

	outputMapper.Apply(inScope, outScope)

	attrs := outScope.GetAttrs()
	inputs := make(map[string]interface{}, len(attrs))

	for name, attr := range attrs {
		inputs[name] = attr.Value
	}

	return inputs
}
