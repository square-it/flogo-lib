package engine

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/service"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// Environment defines the environment in which the engine will run
type Environment struct {
	flowProvider         service.FlowProviderService
	stateRecorder        service.StateRecorderService
	stateRecorderEnabled bool
	engineTester         service.EngineTesterService
	engineTesterEnabled  bool

	engineConfig       *Config
	triggersConfig     *TriggersConfig
	embeddedFlowManger *util.EmbeddedFlowManager
}

// NewEnvironment creates a new engine Environment from the provided configuration and the specified
// StateRecorder and FlowProvider
func NewEnvironment(flowProvider service.FlowProviderService, stateRecorder service.StateRecorderService, engineTester service.EngineTesterService, config *Config, triggersConfig *TriggersConfig) *Environment {

	var engineEnv Environment

	if flowProvider == nil {
		panic("Engine Environment: FlowProvider Service cannot be nil")
	}

	engineEnv.flowProvider = flowProvider
	engineEnv.stateRecorder = stateRecorder
	engineEnv.engineTester = engineTester
	engineEnv.engineConfig = config

	if len(triggersConfig.Triggers) == 0 { //temporary
		engineEnv.triggersConfig.Triggers = config.Triggers
	} else {
		engineEnv.triggersConfig = triggersConfig
	}



	return &engineEnv
}

// FlowProviderService returns the flow.Provider service associated with the EngineEnv
func (e *Environment) FlowProviderService() service.FlowProviderService {
	return e.flowProvider
}

// StateRecorderService returns the flowinst.StateRecorder service associated with the EngineEnv
func (e *Environment) StateRecorderService() (stateRecorder service.StateRecorderService, enabled bool) {

	return e.stateRecorder, e.stateRecorderEnabled
}

// EngineTesterService returns the EngineTester service associated with the EngineEnv
func (e *Environment) EngineTesterService() (engineTester service.EngineTesterService, enabled bool) {

	return e.engineTester, e.engineTesterEnabled
}

// SetEmbeddedJSONFlows sets the embedded flows (in JSON) for the engine
func (e *Environment) SetEmbeddedJSONFlows(compressed bool, jsonFlows map[string]string) {
	e.embeddedFlowManger = util.NewEmbeddedFlowManager(compressed, jsonFlows)
}

// EngineConfig returns the Engine Config for the Engine Environment
func (e *Environment) EngineConfig() *Config {
	return e.engineConfig
}

// TriggersConfig returns the Triggers Config for the Engine Environment
func (e *Environment) TriggersConfig() *TriggersConfig {
	return e.triggersConfig
}

// Init is used to initialize the engine environment
func (e *Environment) Init(instManager *flowinst.Manager, defaultRunner runner.Runner) {

	settings, enabled := getServiceSettings(e.engineConfig, service.ServiceFlowProvider)
	e.flowProvider.Init(settings, e.embeddedFlowManger)

	settings, enabled = getServiceSettings(e.engineConfig, service.ServiceStateRecorder)

	if enabled {
		e.stateRecorderEnabled = true
		e.stateRecorder.Init(settings)
	}

	settings, enabled = getServiceSettings(e.engineConfig, service.ServiceEngineTester)
	if enabled {
		e.engineTesterEnabled = true
		e.engineTester.Init(settings, instManager, defaultRunner)
	}
}

func getServiceSettings(engineConfig *Config, serviceName string) (settings map[string]string, enabled bool) {

	serviceConfig := engineConfig.Services[serviceName]

	enabled = serviceConfig != nil && serviceConfig.Enabled

	if serviceConfig == nil || serviceConfig.Settings == nil {
		settings = make(map[string]string)
	} else {
		settings = serviceConfig.Settings
	}

	return settings, enabled
}
