package engine

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
)

// Environment defines the environment in which the engine will run
type Environment struct {
	processProvider process.Provider
	stateRecorder   processinst.StateRecorder
	engineTester    Tester
	engineConfig    *Config
}

// NewEnvironment creates a new engine Environment from the provided configuration and the specified
// StateRecorder and ProcessProvider
func NewEnvironment(processProvider process.Provider, stateRecorder processinst.StateRecorder, tester Tester, config *Config) *Environment {

	var engineEnv Environment

	engineEnv.processProvider = processProvider
	engineEnv.stateRecorder = stateRecorder
	engineEnv.engineConfig = config
	engineEnv.engineTester = tester

	return &engineEnv
}

// ProcessProvider returns the process.Provider associated with the EngineEnv
func (e *Environment) ProcessProvider() process.Provider {
	return e.processProvider
}

// StateRecorder returns the StateRecorder associated with the EngineEnv
func (e *Environment) StateRecorder() processinst.StateRecorder {
	return e.stateRecorder
}

// EngineTester returns the EngineTester associated with the EngineEnv
func (e *Environment) EngineTester() Tester {
	return e.engineTester
}

// EngineConfig returns the Engine Config for the EngineEnv
func (e *Environment) EngineConfig() *Config {
	return e.engineConfig
}
