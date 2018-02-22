package action

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Action is an action to perform as a result of a trigger
type Action interface {

	//Config get the Action's config
	//Config() *Config

	//Metadata get the Action's metadata
	Metadata() *Metadata

	//IOMetadata get the Action's IO metadata
	IOMetadata() *data.IOMetadata

	// Run this Action
	//Run(context context.Context, inputs []*data.Attribute, options map[string]interface{}) (map[string]*data.Attribute, error)

	// Run this Action
	//Run(context context.Context, inputs []*data.Attribute, options map[string]interface{}, handler ResultHandler) error
}

// SyncAction is a synchronous action to perform as a result of a trigger
type SyncAction interface {
	Action

	// Run this Action
	Run(context context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error)
}

// AsyncAction is an asynchronous action to perform as a result of a trigger, the action can asynchronously
// return results as it runs.  It returns immediately, but will continue to run.
type AsyncAction interface {
	Action

	// Run this Action
	Run(context context.Context, inputs map[string]*data.Attribute, handler ResultHandler) error
}

// Factory is used to create new instances for an action
type Factory interface {

	//New create a new Action
	New(config *Config) (Action, error)
}

// Runner runs actions
type Runner interface {
	//DEPRECATED
	Run(context context.Context, act Action, uri string, options interface{}) (code int, data interface{}, err error)

	//Run the specified Action
	//DEPRECATED
	RunAction(ctx context.Context, act Action, options map[string]interface{}) (results map[string]*data.Attribute, err error)

	//Execute the specified Action
	Execute(ctx context.Context, act Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error)
}

// ResultHandler used to handle results from the Action
type ResultHandler interface {

	// HandleResult invoked when there are results available
	HandleResult(results map[string]*data.Attribute, err error)

	// Done indicates that the action has completed
	Done()
}
