package flowinst

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"github.com/TIBCOSoftware/flogo-lib/core/support"
)

// Starter interface is used to start flow instances, used by Triggers
// to start instances
type Starter interface {

	// StartFlowInstance starts a flow instance using the provided information
	StartFlowInstance(flowURI string, startAttrs []*data.Attribute, replyHandler support.ReplyHandler, execOptions *ExecOptions) (instanceID string, startError error)
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
