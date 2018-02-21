package handler

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"context"
)

type Handler struct {
	metadata *trigger.Metadata
	runner action.Runner
	Action action.Action
	Config *Config
}

func New(config *Config, act action.Action, metadata trigger.Metadata, actionRunner action.Runner) {

}

func (h *Handler) Handle(ctx context.Context, data map[string]interface{}) (results map[string]*data.Attribute, err error){


	//todo handle error
	startAttrs, _ := h.metadata.OutputsToAttrs(data, false)

	newCtx := trigger.NewContext(ctx, startAttrs) //, h.Config)
	results, err = h.runner.RunAction(newCtx, h.Action, nil)

	return results, err
}
