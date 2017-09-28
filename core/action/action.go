package action

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Action is an action to perform as a result of a trigger
type Action interface {

	//Config get the Action's config
	Config() *Config

	//Metadata get the Action's metadata
	Metadata() *Metadata

	// Run this Action
	Run(context context.Context, inputs map[string]interface{}, options map[string]interface{}, handler ResultHandler) error
}

// Factory is used to create new instances for an action
type Factory interface {
	New(config *Config) Action
}

// Runner runs actions
type Runner interface {
	Run(context context.Context, action Action, uri string, options interface{}) (code int, data interface{}, err error)

	//Run the specified Action
	RunAction(context context.Context, actionID string, inputGenerator InputGenerator, options map[string]interface{}) (results map[string]interface{}, err error)
}

//TODO REVIEW - Need a way to package the output of the trigger, its metadata, and the corresponding mapper
//TODO        - so it doesn't get evaluated until Action.Run, probably needs a better name
type InputGenerator interface {
	GenerateInputs(inputMetadata map[string]*data.Attribute) map[string]interface{}
}

// ResultHandler used to handle results from the Action
type ResultHandler interface {

	HandleResult(results map[string]interface{}, err error)

	Done()
}
