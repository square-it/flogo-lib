package tester

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("tester")

// RequestFlowor flowes request objects and invokes the corresponding
// flow Manager methods
type RequestFlowor struct {
	instManager *flowinst.Manager
}

// NewRequestFlowor creates a new Request Flowor
func NewRequestFlowor(instManager *flowinst.Manager) *RequestFlowor {

	var rp RequestFlowor
	rp.instManager = instManager

	return &rp
}

// StartFlow handles a StartRequest for a FlowInstance.  This will
// generate an ID for the new FlowInstance and queue a StartRequest.
func (rp *RequestFlowor) StartFlow(startRequest *StartRequest, replyHandler flowinst.ReplyHandler) *flowinst.Instance {

	execOptions := &flowinst.ExecOptions{Interceptor: startRequest.Interceptor, Patch: startRequest.Patch}
	instance := rp.instManager.StartInstance(startRequest.FlowURI, startRequest.Data, replyHandler, execOptions)

	return instance
}

// RestartFlow handles a RestartRequest for a FlowInstance.  This will
// generate an ID for the new FlowInstance and queue a RestartRequest.
func (rp *RequestFlowor) RestartFlow(restartRequest *RestartRequest, replyHandler flowinst.ReplyHandler) *flowinst.Instance {

	execOptions := &flowinst.ExecOptions{Interceptor: restartRequest.Interceptor, Patch: restartRequest.Patch}
	instance := rp.instManager.RestartInstance(restartRequest.IntialState, restartRequest.Data, replyHandler, execOptions)

	return instance
}

// ResumeFlow handles a ResumeRequest for a FlowInstance.  This will
// queue a RestartRequest.
func (rp *RequestFlowor) ResumeFlow(resumeRequest *ResumeRequest, replyHandler flowinst.ReplyHandler) *flowinst.Instance {

	execOptions := &flowinst.ExecOptions{Interceptor: resumeRequest.Interceptor, Patch: resumeRequest.Patch}
	instance := rp.instManager.ResumeInstance(resumeRequest.State, resumeRequest.Data, replyHandler, execOptions)

	return instance
}

// StartRequest describes a request for starting a FlowInstance
type StartRequest struct {
	FlowURI     string            `json:"flowUri"`
	Data        map[string]string `json:"data"`
	Interceptor *flow.Interceptor `json:"interceptor"`
	Patch       *flow.Patch       `json:"patch"`
	ReplyTo     string            `json:"replyTo"`
}

// RestartRequest describes a request for restarting a FlowInstance
// todo: can be merged into StartRequest
type RestartRequest struct {
	IntialState *flowinst.Instance `json:"initialState"`
	Data        map[string]string  `json:"data"`
	Interceptor *flow.Interceptor  `json:"interceptor"`
	Patch       *flow.Patch        `json:"patch"`
}

// ResumeRequest describes a request for resuming a FlowInstance
//todo: Data for resume request should be directed to wating task
type ResumeRequest struct {
	State       *flowinst.Instance `json:"state"`
	Data        map[string]string  `json:"data"`
	Interceptor *flow.Interceptor  `json:"interceptor"`
	Patch       *flow.Patch        `json:"patch"`
}
