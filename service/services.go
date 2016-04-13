package service

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

const (
	// ServiceStateRecorder is the name of the StateRecorder service used in configuration
	ServiceStateRecorder string = "stateRecorder"

	// ServiceProcessProvider is the name of the ProcessProvider service used in configuration
	ServiceProcessProvider string = "processProvider"

	// ServiceEngineTester is the name of the EngineTester service used in configuration
	ServiceEngineTester string = "engineTester"
)

// StateRecorderService is the processinst.StateRecorder wrapped as a service
type StateRecorderService interface {
	util.Managed
	processinst.StateRecorder

	Init(settings map[string]string)
}

// ProcessProviderService is the process.Provider wrapped as a service
type ProcessProviderService interface {
	util.Managed
	process.Provider

	Init(settings map[string]string, embeddedFlowMgr *util.EmbeddedFlowManager)
}

// EngineTesterService is an engine service to assist in testing processes
type EngineTesterService interface {
	util.Managed

	//Init initializes the EngineTester
	Init(settings map[string]string, instManager *processinst.Manager, runner runner.Runner)
}
