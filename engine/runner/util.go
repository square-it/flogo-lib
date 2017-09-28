package runner

import (
	"context"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

type OldTAInputGenerator struct {
	ctx context.Context
}

func NewOldTAInputGenerator(ctx context.Context) *OldTAInputGenerator {

	return &OldTAInputGenerator{ctx: ctx}
}

func (ig *OldTAInputGenerator) GenerateInputs(inputMetadata map[string]*data.Attribute) map[string]interface{} {

	if ig.ctx == nil {
		return nil
	}

	triggerAttrs, ok := trigger.FromContext(ig.ctx)

	if ok {
		attrs := make(map[string]interface{})

		if len(triggerAttrs) > 0 {
			logger.Debug("Run Attributes:")

			for _, attr := range triggerAttrs {

				logger.Debugf(" Attr: %s, Type: %s, Value: %v", attr.Name, attr.Type.String(), attr.Value)
				attrName := "_T." + attr.Name

				attrs[attrName] = data.NewAttribute(attrName, attr.Type, attr.Value)

				//// Keep Temporarily, for short ttterm backwards compatibility
				//attrName1 := "{T." + attr.Name + "}"
				//
				//attrName2 := "{TriggerData." + attr.Name + "}"
				//attrs[attrName2] = data.NewAttribute(attrName2, attr.Type, attr.Value)
				//
				//attrName3 := "${trigger." + attr.Name + "}"
				//attrs[attrName3] = data.NewAttribute(attrName3, attr.Type, attr.Value)
			}
		}

		return attrs
	}

	return nil
}
