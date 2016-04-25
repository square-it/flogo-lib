package flowinst

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
)

// Starter interface is used to start flow instances, used by Triggers
// to start instances
type Starter interface {

	// StartFlowInstance starts a flow instance using the provided information
	StartFlowInstance(flowURI string, startAttrs []*data.Attribute, replyHandler ReplyHandler, execOptions *ExecOptions) (instanceID string, startError error)
}

// ReplyHandler is used to reply back to whoever started the flow instance
type ReplyHandler interface {

	// Reply is used to reply with the results of the instance execution
	Reply(replyData map[string]string)
}

// ExecOptions are optional Patch & Interceptor to be used during instance execution
type ExecOptions struct {
	Patch       *flow.Patch
	Interceptor *flow.Interceptor
}

// IDGenerator generates IDs for flow instances
type IDGenerator interface {

	//NewFlowInstanceID generate a new instance ID
	NewFlowInstanceID() string
}
