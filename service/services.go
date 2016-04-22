package service

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

const (
	// ServiceStateRecorder is the name of the StateRecorder service used in configuration
	ServiceStateRecorder string = "stateRecorder"

	// ServiceFlowProvider is the name of the FlowProvider service used in configuration
	ServiceFlowProvider string = "flowProvider"

	// ServiceEngineTester is the name of the EngineTester service used in configuration
	ServiceEngineTester string = "engineTester"
)

// StateRecorderService is the flowinst.StateRecorder wrapped as a service
type StateRecorderService interface {
	util.Managed
	flowinst.StateRecorder

	Init(settings map[string]string)
}

// FlowProviderService is the flow.Provider wrapped as a service
type FlowProviderService interface {
	util.Managed
	flow.Provider

	Init(settings map[string]string, embeddedFlowMgr *util.EmbeddedFlowManager)
}

// EngineTesterService is an engine service to assist in testing flowes
type EngineTesterService interface {
	util.Managed

	//Init initializes the EngineTester
	Init(settings map[string]string)

	//SetupInstanceSupport sets up the support for instance execution
	SetupInstanceSupport(instManager *flowinst.Manager, runner runner.Runner)
}
