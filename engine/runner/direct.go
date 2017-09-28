package runner

import (
	"context"
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
)

// DirectRunner runs an action synchronously
type DirectRunner struct {
}

// NewDirectRunner create a new DirectRunner
func NewDirect() *DirectRunner {
	return &DirectRunner{}
}

// Start will start the engine, by starting all of its workers
func (runner *DirectRunner) Start() error {
	//op-op
	return nil
}

// Stop will stop the engine, by stopping all of its workers
func (runner *DirectRunner) Stop() error {
	//no-op
	return nil
}

//Run
//Deprecated
func (runner *DirectRunner) Run(ctx context.Context, act action.Action, uri string, options interface{}) (code int, data interface{}, err error) {

	if act == nil {
		return 0, nil, errors.New("Action not found")
	}

	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = options

	handler := &SyncResultHandler{done: make(chan bool, 1)}

	inputGenerator := NewOldTAInputGenerator(ctx)
	inputs := inputGenerator.GenerateInputs(action.GetConfigInputMetadata(act))

	err = act.Run(ctx, inputs, newOptions, handler)

	if err != nil {
		return 0, nil, err
	}

	<-handler.done

	ndata, err  := handler.Result()

	//ndata, err := runner.RunAction(ctx, uri, NewOldTAInputGenerator(ctx), newOptions)

	if len(ndata) != 0 {
		defData, ok := ndata["data"]
		if ok {
			data = defData
		}
		defCode, ok := ndata["code"]
		if ok {
			code = defCode.(int)
		}
	}

	return code, data, err
}

// Run the specified action
func (runner *DirectRunner) RunAction(ctx context.Context, actionID string, inputGenerator action.InputGenerator, options map[string]interface{}) (results map[string]interface{}, err error) {

	act := action.Get(actionID)

	if act == nil {
		return nil, errors.New("Action not found")
	}

	handler := &SyncResultHandler{done: make(chan bool, 1)}

	inputs := inputGenerator.GenerateInputs(action.GetConfigInputMetadata(act))

	err = act.Run(ctx, inputs, options, handler)

	if err != nil {
		return nil, err
	}

	<-handler.done

	return handler.Result()
}

// SyncResultHandler simple result handler to use in synchronous case
type SyncResultHandler struct {
	done chan (bool)
	data map[string]interface{}
	err  error
}

// HandleResult implements action.ResultHandler.HandleResult
func (rh *SyncResultHandler) HandleResult(data map[string]interface{}, err error) {
	rh.data = data
	rh.err = err
}

// Done implements action.ResultHandler.Done
func (rh *SyncResultHandler) Done() {
	rh.done <- true
}

// Result returns the latest Result set on the handler
func (rh *SyncResultHandler) Result() (data map[string]interface{}, err error) {
	return rh.data, rh.err
}
