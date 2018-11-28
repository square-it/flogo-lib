package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"io/ioutil"
)

func init() {
	profiles := config.GetAppConfigProfiles()
	if profiles != "" {
		file, e := ioutil.ReadFile(profiles)
		if e != nil {
			return // TODO : log
		}

		properties := make(map[string]interface{})
		e = json.Unmarshal(file, &properties)

		if e != nil {
			return // TODO : log
		}

		RegisterPropertyValueResolver("file", &FileValueResolver{properties: properties})
	}
}

// Resolve property value from a JSON profile file
type FileValueResolver struct {
	properties map[string]interface{}
}

func (resolver *FileValueResolver) ResolveValue(key string) (interface{}, error) {
	value, exists := resolver.properties[key]
	if !exists {
		return nil, errors.New(fmt.Sprintf("Environment variable '%s' is not set", key))
	}
	return value, nil
}
