package trigger

import (
	"context"
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var handlerMap map[*HandlerConfig]*Handler

//DEPRECATED
type LegacyRunner struct {
	currentRunner action.Runner
}

func NewLegacyRunner(runner action.Runner) action.Runner {
	return &LegacyRunner{currentRunner: runner}
}

//RegisterHandler register handler for the specified configuration
//DEPRECATED
func RegisterHandler(cfg *HandlerConfig, h *Handler) {
	if handlerMap == nil {
		handlerMap = make(map[*HandlerConfig]*Handler)
	}
	handlerMap[cfg] = h
}

func (lr *LegacyRunner) Run(ctx context.Context, act action.Action, uri string, options interface{}) (code int, data interface{}, err error) {

	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = options

	results, err := lr.RunAction(ctx, act, newOptions)

	if len(results) != 0 {
		defData, ok := results["data"]
		if ok {
			data = defData.Value()
		}
		defCode, ok := results["code"]
		if ok {
			code = defCode.Value().(int)
		}
	}

	return code, data, err
}

func (lr *LegacyRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {

	trgHandler, trgData := lr.getHandler(ctx, act)
	return trgHandler.Handle(ctx, trgData)
}

func (*LegacyRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	//only called by handler so not needed

	return nil, errors.New("not supported")
}

func (lr *LegacyRunner) getHandler(ctx context.Context, act action.Action) (*Handler, map[string]interface{}) {
	var values map[string]interface{}
	var handler *Handler

	if ctx != nil {
		var exists bool
		ctxData, exists := ExtractContextData(ctx)

		if !exists {
			values = attrsToData(ctxData.Attrs)
			handler = handlerMap[ctxData.HandlerCfg]
		}
	}

	if handler == nil {
		logger.Warn("unable to find existing handler, creating new one")
		handler = NewHandler(nil, act, nil, nil, lr.currentRunner)
	}

	return handler, values
}

func attrsToData(attrs []*data.Attribute) map[string]interface{} {

	if attrs == nil {
		return nil
	}

	values := make(map[string]interface{}, len(attrs))

	for _, attr := range attrs {
		values[attr.Name()] = attr
	}

	return values
}
