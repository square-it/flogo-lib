package engine

import "github.com/TIBCOSoftware/flogo-lib/service"


// Environment defines the environment in which the engine will run
type Environment struct {
	processProvider service.ProcessProviderService
	stateRecorder   service.StateRecorderService
	engineTester    service.EngineTesterService
	engineConfig    *Config
}

// NewEnvironment creates a new engine Environment from the provided configuration and the specified
// StateRecorder and ProcessProvider
func NewEnvironment(processProvider service.ProcessProviderService, stateRecorder service.StateRecorderService, engineTester service.EngineTesterService, config *Config) *Environment {

	var engineEnv Environment

	if processProvider == nil {
		panic("Engine Environment: ProcessProvider Service cannot be nil")
	}

	engineEnv.processProvider = processProvider
	engineEnv.stateRecorder = stateRecorder
	engineEnv.engineTester = engineTester
	engineEnv.engineConfig = config

	return &engineEnv
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) ProcessProviderService() service.ProcessProviderService {
	return e.processProvider
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) ProcessProviderServiceSettings() (settings map[string]string, enabled bool) {
	settings, enabled =  getServiceSettings(e.engineConfig, service.ServiceProcessProvider)
	return settings, enabled && e.processProvider != nil
}

// StateRecorderService returns the StateRecorder service associated with the EngineEnv
func (e *Environment) StateRecorderService() service.StateRecorderService {
	return e.stateRecorder
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) StateRecorderServiceSettings() (settings map[string]string, enabled bool) {
	settings, enabled = getServiceSettings(e.engineConfig, service.ServiceStateRecorder)
	return settings, enabled && e.stateRecorder != nil
}

// EngineTesterService returns the EngineTester service associated with the EngineEnv
func (e *Environment) EngineTesterService() service.EngineTesterService {
	return e.engineTester
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) EngineTesterServiceSettings() (settings map[string]string, enabled bool) {
	settings, enabled =  getServiceSettings(e.engineConfig, service.ServiceEngineTester)
	return settings, enabled && e.engineTester != nil
}

// EngineConfig returns the Engine Config for the Engine Environment
func (e *Environment) EngineConfig() *Config {
	return e.engineConfig
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

