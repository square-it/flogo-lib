package starter

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
)

//todo consolidate start & restart
type ProcessStarter interface {

	// StartProcess handles a StartRequest for a ProcessInstance.  This will
	// generate an ID for the new ProcessInstance and queue a StartRequest.
	StartProcess(startRequest *StartRequest) string

	// RestartProcess handles a RestartRequest for a ProcessInstance.  This will
	// generate an ID for the new ProcessInstance and queue a RestartRequest.
	RestartProcess(restartRequest *RestartRequest) string

	// ResumeProcess handles a ResumeRequest for a ProcessInstance.  This will
	// queue a RestartRequest.
	ResumeProcess(resumeRequest *ResumeRequest)
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string               `json:"processUri"`
	Data        map[string]string    `json:"data"`
	Interceptor *process.Interceptor `json:"interceptor"`
	Patch       *process.Patch       `json:"patch"`
	ReplyTo	    string               `json:"replyTo"`
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

