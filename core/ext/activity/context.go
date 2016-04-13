package activity

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Context describes the execution context for an Activity.
// It provides access to attributes, task and Flow information.
type Context interface {
	data.Scope

	// FlowInstanceID returns the ID of the Flow Instance
	FlowInstanceID() string

	// FlowName returns the name of the Flow
	FlowName() string

	// TaskName returns the name of the Task the Activity is currently executing
	TaskName() string
}
