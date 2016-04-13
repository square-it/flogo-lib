package trigger

// Config is the configuration for a Trigger
// todo: should we switch to settings to map[string]interface{} ?
type Config struct {
	Name      string            `json:"name"`
	Settings  map[string]string `json:"settings"`
	Endpoints []*EndpointConfig `json:"endpoints"`
}

// EndpointConfig is the configuration for a specific endpoint for the
// Trigger
type EndpointConfig struct {
	FlowURI  string            `json:"flowURI"`
	Settings map[string]string `json:"settings"`
}
