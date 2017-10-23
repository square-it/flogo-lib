package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)


/*
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

 */
func generateInputs(act action.Action, ctxData *trigger.ContextData) ([]*data.Attribute) {

	if ctxData == nil || ctxData.Attrs == nil {
		return nil
	}

	inputMetadata := action.GetConfigInputMetadata(act)

	if ctxData.HandlerCfg != nil && inputMetadata != nil {

		outputMapper := ctxData.HandlerCfg.GetActionInputMapper()

		outScope := data.NewFixedScope(inputMetadata)
		inScope := data.NewSimpleScope(ctxData.Attrs, nil)

		outputMapper.Apply(inScope, outScope)

		attrs := outScope.GetAttrs()

		inputs := make([]*data.Attribute, 0, len(inputMetadata))

		for _, attr := range attrs {
			inputs = append(inputs, attr)
		}

		return inputs

	} else {
		// for backwards compatibility make trigger outputs map directly to action inputs

		if len(ctxData.Attrs) > 0 {
			logger.Debug("No mapping specified, adding trigger outputs as inputs to action")

			inputs := make([]*data.Attribute, 0, len(ctxData.Attrs))

			for _, attr := range ctxData.Attrs {

				logger.Debugf(" Attr: %s, Type: %s, Value: %v", attr.Name, attr.Type.String(), attr.Value)
				attrName := "_T." + attr.Name

				inputs = append(inputs, data.NewAttribute(attrName, attr.Type, attr.Value))
			}

			return inputs
		}
	}

	return nil
}

