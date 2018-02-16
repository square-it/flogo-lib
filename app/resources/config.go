package resources

import "encoding/json"

type ResourcesConfig struct {
	Resources []*ResourceConfig `json:"resources"`
}

type ResourceConfig struct {
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}
