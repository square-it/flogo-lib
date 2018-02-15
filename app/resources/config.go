package resources

import "encoding/json"

type ResourcesConfig struct {
	Resources []*ResourceSetConfig `json:"resources"`
}

type ResourceSetConfig struct {
	Type string `json:"type"`
	Entries []json.RawMessage `json:"entries"`
}