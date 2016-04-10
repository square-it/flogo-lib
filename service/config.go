package service


type Config struct {
	Name     string            `json:"name"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]string `json:"settings,omitempty"`
}

func DefaultServicesConfig() map[string]*Config {
	servicesCfg := make(map[string]*Config)

	//todo: move to individual service implementations or probably remove 'default' implemenations
	servicesCfg[ServiceStateRecorder] = &Config{Name:ServiceStateRecorder, Enabled: true, Settings: map[string]string{"host": ""}}
	servicesCfg[ServiceProcessProvider] = &Config{Name:ServiceProcessProvider, Enabled: true}
	servicesCfg[ServiceEngineTester] = &Config{Name: ServiceEngineTester, Enabled: true, Settings: map[string]string{"port": "8080"}}

	return servicesCfg
}