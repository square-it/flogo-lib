package engine

import "github.com/TIBCOSoftware/flogo-lib/services"


// Environment defines the environment in which the engine will run
type Environment struct {
	providerService     services.ProcessProviderService
	recorderService     services.StateRecorderService
	engineTesterService services.TesterService
	engineConfig        *Config
}

// NewEnvironment creates a new engine Environment from the provided configuration and the specified
// StateRecorder and ProcessProvider
func NewEnvironment(providerService services.ProcessProviderService, recorderService services.StateRecorderService, testerService services.TesterService, config *Config) *Environment {

	var engineEnv Environment

	if providerService == nil {
		panic("Engine Environment: ProcessProvider Service cannot be nil")
	}

	engineEnv.providerService = providerService
	engineEnv.recorderService = recorderService
	engineEnv.engineTesterService = testerService
	engineEnv.engineConfig = config

	return &engineEnv
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) ProcessProviderService() services.ProcessProviderService {
	return e.providerService
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) ProcessProviderServiceSettings() (settings map[string]string, enabled bool) {
	settings, enabled =  getServiceSettings(e.engineConfig, services.ServiceProcessProvider)
	return settings, enabled && e.providerService != nil
}

// StateRecorderService returns the StateRecorder service associated with the EngineEnv
func (e *Environment) StateRecorderService() services.StateRecorderService {
	return e.recorderService
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) StateRecorderServiceSettings() (settings map[string]string, enabled bool) {
	settings, enabled = getServiceSettings(e.engineConfig, services.ServiceStateRecorder)
	return settings, enabled && e.recorderService != nil
}

// EngineTesterService returns the EngineTester service associated with the EngineEnv
func (e *Environment) EngineTesterService() services.TesterService {
	return e.engineTesterService
}

// ProcessProviderService returns the process.Provider service associated with the EngineEnv
func (e *Environment) EngineTesterServiceSettings() (settings map[string]string, enabled bool) {
	settings, enabled =  getServiceSettings(e.engineConfig, services.ServiceEngineTester)
	return settings, enabled && e.engineTesterService != nil
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

