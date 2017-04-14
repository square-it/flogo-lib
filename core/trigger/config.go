package trigger

// Config is the configuration for a Trigger
type Config struct {
	Name     string            `json:"name"`
	Id       string            `json:"id"`
	Ref      string            `json:"ref"`
	Settings map[string]string `json:"settings"`
	Handlers []*HandlerConfig  `json:"handlers"`

	//deprecated
	//Settings map[string]string `json:"settings"`
	Endpoints []*EndpointConfig `json:"endpoints"`
}

// HandlerConfig is the configuration for the Trigger Handler
type HandlerConfig struct {
	ActionId string            `json:"actionId"`
	Settings map[string]string `json:"settings"`
}


//// Trigger is the configuration for the Trigger
//type TriggerConfig struct {
//	Id       string                 `json:"id"`
//	Ref      string                 `json:"ref"`
//	Settings map[string]interface{} `json:"settings"`
//	Handlers []*TriggerHandler      `json:"handlers"`
//}

// EndpointConfig is the configuration for a specific endpoint for the
// Trigger // Deprecated
type EndpointConfig struct {
	ActionId   string            `json:"actionId"`
	ActionType string            `json:"actionType"`
	ActionURI  string            `json:"actionURI"`
	Settings   map[string]string `json:"settings"`
}

