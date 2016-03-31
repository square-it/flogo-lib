package activity

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Context describes the execution context for an Activity.
// It provides access to attributes, task and Process information.
type Context interface {
	data.Scope

	// ProcessInstanceID returns the ID of the Process Instance
	ProcessInstanceID() string

	// ProcessName returns the name of the Process
	ProcessName() string

	// TaskName returns the name of the Task the Activity is currently executing
	TaskName() string
}
