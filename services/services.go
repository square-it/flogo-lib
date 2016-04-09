package services

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("services")

const (
	ServiceStateRecorder string = "stateRecorder"
	ServiceProcessProvider string = "processProvider"
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

	Init(settings map[string]string)
}

// TesterService is an engine service to assist in testing processes
type TesterService interface {
	util.Managed

	//Init initializes the EngineTester
	Init(instManager *processinst.Manager, runner runner.Runner, settings map[string]string)
}

