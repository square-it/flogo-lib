package engine

import (
	"github.com/op/go-logging"
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
)

var log = logging.MustGetLogger("engine")

// System is core configuration, managment and environment object, which
// is central to an Engine
type System struct {
	processProvider process.Provider
	stateRecorder   processinst.StateRecorder
	engine         *Engine
}

// NewSystem creates a new system from the provided configuration and the specified
// StateRecorder and ProcessProvider
func NewSystem(processProvider process.Provider, stateRecorder processinst.StateRecorder, config *EngineConfig) *System {

	var system System

	system.processProvider = processProvider
	system.stateRecorder = stateRecorder
	system.engine = NewEngine(&system, config)
	return &system
}

// ProcessProvider returns the process.Provider associated with the System
func (sys *System) ProcessProvider() process.Provider {
	return sys.processProvider
}

// StateRecorder returns the StateRecorder associated with the System
func (sys *System) StateRecorder() processinst.StateRecorder {
	return sys.stateRecorder
}

// GetEngine returns the Engine for the System
func (sys *System) GetEngine() *Engine {
	return sys.engine
}
