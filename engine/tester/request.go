package tester

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("tester")

// RequestProcessor processes request objects and invokes the corresponding
// process Manager methods
type RequestProcessor struct {
	instManager *processinst.Manager
}

// NewRequestProcessor creates a new Request Processor
func NewRequestProcessor(instManager *processinst.Manager) *RequestProcessor {

	var rp RequestProcessor
	rp.instManager = instManager

	return &rp
}

// StartProcess handles a StartRequest for a ProcessInstance.  This will
// generate an ID for the new ProcessInstance and queue a StartRequest.
func (rp *RequestProcessor) StartProcess(startRequest *StartRequest, replyHandler processinst.ReplyHandler) *processinst.Instance {

	execOptions := &processinst.ExecOptions{Interceptor: startRequest.Interceptor, Patch: startRequest.Patch}
	instance := rp.instManager.StartInstance(startRequest.ProcessURI, startRequest.Data, replyHandler, execOptions)

	return instance
}

// RestartProcess handles a RestartRequest for a ProcessInstance.  This will
// generate an ID for the new ProcessInstance and queue a RestartRequest.
func (rp *RequestProcessor) RestartProcess(restartRequest *RestartRequest, replyHandler processinst.ReplyHandler) *processinst.Instance {

	execOptions := &processinst.ExecOptions{Interceptor: restartRequest.Interceptor, Patch: restartRequest.Patch}
	instance := rp.instManager.RestartInstance(restartRequest.IntialState, restartRequest.Data, replyHandler, execOptions)

	return instance
}

// ResumeProcess handles a ResumeRequest for a ProcessInstance.  This will
// queue a RestartRequest.
func (rp *RequestProcessor) ResumeProcess(resumeRequest *ResumeRequest, replyHandler processinst.ReplyHandler) *processinst.Instance {

	execOptions := &processinst.ExecOptions{Interceptor: resumeRequest.Interceptor, Patch: resumeRequest.Patch}
	instance := rp.instManager.ResumeInstance(resumeRequest.State, resumeRequest.Data, replyHandler, execOptions)

	return instance
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string               `json:"processUri"`
	Data        map[string]string    `json:"data"`
	Interceptor *process.Interceptor `json:"interceptor"`
	Patch       *process.Patch       `json:"patch"`
	ReplyTo     string               `json:"replyTo"`
}

// RestartRequest describes a request for restarting a ProcessInstance
// todo: can be merged into StartRequest
type RestartRequest struct {
	IntialState *processinst.Instance `json:"initialState"`
	Data        map[string]string     `json:"data"`
	Interceptor *process.Interceptor  `json:"interceptor"`
	Patch       *process.Patch        `json:"patch"`
}

// ResumeRequest describes a request for resuming a ProcessInstance
//todo: Data for resume request should be directed to wating task
type ResumeRequest struct {
	State       *processinst.Instance `json:"state"`
	Data        map[string]string     `json:"data"`
	Interceptor *process.Interceptor  `json:"interceptor"`
	Patch       *process.Patch        `json:"patch"`
}
