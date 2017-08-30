package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
)

type OldTAInputGenerator struct {
	ctx context.Context
}

func NewOldTAInputGenerator(ctx context.Context) *OldTAInputGenerator {

	return &OldTAInputGenerator{ctx: ctx}
}

func (ig *OldTAInputGenerator) GenerateInputs(inputMetadata map[string]*data.Attribute) map[string]interface{} {

	triggerAttrs, ok := trigger.FromContext(ig.ctx)

	if ok {
		attrs := make(map[string]interface{})

		if len(triggerAttrs) > 0 {
			logger.Debug("Run Attributes:")

			for _, attr := range triggerAttrs {

				logger.Debugf(" Attr:%s, Type:%s, Value:%v", attr.Name, attr.Type.String(), attr.Value)

				// Keep Temporarily, for short term backwards compatibility
				attrName1 := "{T." + attr.Name + "}"
				attrs[attrName1] = attr.Value

				attrName2 := "{TriggerData." + attr.Name + "}"
				attrs[attrName2] = attr.Value

				attrName3 := "${trigger." + attr.Name + "}"
				attrs[attrName3] = attr.Value
			}
		}

		return attrs
	}

	return nil
}

func GetActionOutputMetadata(act action.Action) map[string]*data.Attribute {

	if act.Config() != nil {
		if act.Config().Metadata != nil {
			return act.Config().Metadata.Outputs
		}
	}

	return nil
}

