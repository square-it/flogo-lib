package processinst

import "github.com/TIBCOSoftware/flogo-lib/core/process"

// Starter interface is used to start process instances, used by Triggers
// to start instances
type Starter interface {

	// StartProcessInstance starts a process instance using the provided information
	// todo: make data map[string]interface{}
	StartProcessInstance(processURI string, startData map[string]string, replyHandler ReplyHandler, execOptions *ExecOptions) string
}

// ReplyHandler is used to reply back to whoever started the process instance
type ReplyHandler interface {

	// Reply is used to reply with the results of the instance execution
	Reply(replyData map[string]string)
}

// ExecOptions are optional Patch & Interceptor to be used during instance execution
type ExecOptions struct {
	Patch       *process.Patch
	Interceptor *process.Interceptor
}

// IDGenerator generates IDs for process instances
type IDGenerator interface {

	//NewProcessInstanceID generate a new instance ID
	NewProcessInstanceID() string
}
