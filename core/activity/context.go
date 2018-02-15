package activity

import (
	//"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Context describes the execution context for an Activity.
// It provides access to attributes, task and Flow information.
type Context interface {

	// FlowDetails gets the action.Context under with the activity is executing
	ActivityHost() Host

	// TaskName returns the name of the Task the Activity is currently executing
	//Deprecated
	TaskName() string

	//Name the name of the activity that is currently executing
	Name() string

	// GetInput gets the value of the specified input attribute
	GetInput(name string) interface{}

	// GetOutput gets the value of the specified output attribute
	GetOutput(name string) interface{}

	// SetOutput sets the value of the specified output attribute
	SetOutput(name string, value interface{})

	//Deprecated
	// FlowDetails returns the details fo the Flow Instance
	//FlowDetails() FlowDetails
}

// Deprecated
// FlowDetails details of the flow that is being executed
//type FlowDetails interface {
//
//	// ID returns the ID of the Flow Instance
//	ID() string
//
//	// FlowName returns the name of the Flow
//	Name() string
//
//	// ReplyHandler returns the reply handler for the flow Instance
//	ReplyHandler() ReplyHandler
//}

type Host interface {
	// ID returns the ID of the Action Instance
	ID() string

	Name() string

	// The action reference
	//Ref() string

	// IOMetadata get the input/output metadata of the activity host
	IOMetadata() *data.IOMetadata

	// Reply is used to reply to the activity Host with the results of the execution
	Reply(replyData map[string]*data.Attribute, err error)

	// Return is used to indicate to the activity Host that it should complete and return the results of the execution
	Return(returnData map[string]*data.Attribute, err error)

	//todo rename, essentially the flow's attrs for now
	WorkingData() data.Scope

	//Map with action specific details/properties, flowId, etc.
	//GetDetails() map[string]string

	GetResolver() data.Resolver
}