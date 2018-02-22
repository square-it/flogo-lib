package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"context"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
)

type Handler struct {
	runner action.Runner
	act    action.Action

	outputMd map[string]*data.Attribute
	replyMd  map[string]*data.Attribute

	config *HandlerConfig

	actionInputMapper  data.Mapper
	actionOutputMapper data.Mapper
}

func NewHandler(config *HandlerConfig, act action.Action, outputMd map[string]*data.Attribute, replyMd map[string]*data.Attribute, runner action.Runner) *Handler {
	handler := &Handler{config: config, act: act, outputMd: outputMd, replyMd: replyMd, runner: runner}

	if config.Action.Mappings != nil {
		if len(config.Action.Mappings.Input) > 0 {
			handler.actionInputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: config.Action.Mappings.Input}, nil)
		}
		if len(config.Action.Mappings.Output) > 0 {
			handler.actionOutputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: config.Action.Mappings.Output}, nil)
		}
	}

	return handler
}

func (h *Handler) GetSetting(setting string) (interface{}, bool) {
	val, exists := data.GetValueWithResolver(h.config.Settings, setting)

	if !exists {
		val, exists = data.GetValueWithResolver(h.config.parent.Settings, setting)
	}

	return val, exists
}

func (h *Handler) GetStringSetting(setting string) string {
	val, exists := h.GetSetting(setting)

	if !exists {
		return ""
	}

	strVal, err := data.CoerceToString(val)

	if err != nil {
		return ""
	}

	return strVal
}



func (h *Handler) Handle(ctx context.Context, triggerData map[string]interface{}) (map[string]*data.Attribute, error) {

	inputs, err := h.generateInputs(triggerData)

	if err != nil {
		return nil, err
	}

	results, err := h.runner.RunAction2(ctx, h.act, inputs)

	if err != nil {
		return nil, err
	}

	retValue, err := h.generateOutputs(results)

	return retValue, err
}

func (h *Handler) dataToAttrs(triggerData map[string]interface{}) ([]*data.Attribute, error) {

	attrs := make([]*data.Attribute, 0, len(h.outputMd))

	for k, a := range h.outputMd {
		v, _ := triggerData[k]

		var err error
		attr, err := data.NewAttribute(a.Name(), a.Type(), v)
		attrs = append(attrs, attr)

		if err != nil {
			return nil, err
		}
	}

	return attrs, nil
}

func (h *Handler) generateInputs(triggerData map[string]interface{}) ([]*data.Attribute, error) {

	if len(triggerData) == 0 {
		return nil, nil
	}

	inputMetadata := action.GetConfigInputMetadata(h.act)
	triggerAttrs, _ := h.dataToAttrs(triggerData)

	if len(triggerAttrs) == 0 {
		return nil, nil
	}

	var inputs []*data.Attribute

	if h.actionInputMapper != nil && inputMetadata != nil {

		inScope := data.NewSimpleScope(triggerAttrs, nil)
		outScope := data.NewFixedScope(inputMetadata)

		err := h.actionInputMapper.Apply(inScope, outScope)
		if err != nil {
			return nil, err
		}

		attrs := outScope.GetAttrs()

		inputs = make([]*data.Attribute, 0, len(inputMetadata))

		for _, attr := range attrs {
			inputs = append(inputs, attr)
		}
	} else {
		// for backwards compatibility make trigger outputs map directly to action inputs

		logger.Debug("No mapping specified, adding trigger outputs as inputs to action")

		inputs := make([]*data.Attribute, 0, len(triggerAttrs))

		for _, attr := range triggerAttrs {

			logger.Debugf(" Attr: %s, Type: %s, Value: %v", attr.Name(), attr.Type().String(), attr.Value())
			//inputs = append(inputs, data.NewAttribute( attr.Name, attr.Type, attr.Value))
			inputs = append(inputs, attr)

			attrName := "_T." + attr.Name()

			inputs = append(inputs, data.CloneAttribute(attrName, attr))
		}
	}

	return inputs, nil
}

func (h *Handler) generateOutputs(actionResults map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	if len(actionResults) == 0 {
		return nil, nil
	}

	if h.actionOutputMapper == nil {
		//for backwards compatibility
		return actionResults, nil
	}

	outputMetadata := action.GetConfigOutputMetadata(h.act)

	if outputMetadata != nil {

		outScope := data.NewFixedScopeFromMap(h.replyMd)
		inScope := data.NewSimpleScopeFromMap(actionResults, nil)

		err := h.actionOutputMapper.Apply(inScope, outScope)
		if err != nil {
			return nil, err
		}

		return outScope.GetAttrs(), nil
	}

	return actionResults, nil
}
