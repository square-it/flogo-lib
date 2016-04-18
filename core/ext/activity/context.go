package activity

// Context describes the execution context for an Activity.
// It provides access to attributes, task and Flow information.
type Context interface {

	// FlowInstanceID returns the ID of the Flow Instance
	FlowInstanceID() string

	// FlowName returns the name of the Flow
	FlowName() string

	// TaskName returns the name of the Task the Activity is currently executing
	TaskName() string

	// GetInput gets the value of the specified input attribute
	GetInput(name string) interface{}

	// SetOutput sets the value of the specified output attribute
	SetOutput(name string, value interface{})
}
