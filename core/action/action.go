package action

import (
	"context"

	"github.com/TIBCOSoftware/flogo-lib/types"
)

// Action is an action to perform as a result of a trigger
type Action interface {
	// Run this Action
	Run(context context.Context, uri string, options interface{}, handler ResultHandler) error
}

// Action is an action to perform as a result of a trigger
type Action2 interface {
	// Run this Action
	Run(context context.Context, uri string, options interface{}, handler ResultHandler) error

	// Init sets up the action
	Init(config types.ActionConfig)

	// New is a factory function to create a new instance for an id
	New(id string) Action2
}

// Runner runs actions
type Runner interface {
	//Run the specified Action
	Run(context context.Context, action Action, uri string, options interface{}) (code int, data interface{}, err error)
}

// ResultHandler used to handle results from the Action
type ResultHandler interface {
	HandleResult(code int, data interface{}, err error)

	Done()
}
