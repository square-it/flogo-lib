package app

import (
	"encoding/json"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/pkg/errors"
)

// App is the configuration for the App
type Config struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Version     string             `json:"version"`
	Description string             `json:"description"`
	Properties  interface{}        `json:"properties"`
	Triggers    []*trigger.Config  `json:"triggers"`
	Resources   []*resource.Config `json:"resources"`

	//for backwards compatibility
	Actions []*action.Config `json:"actions"`
}

// defaultConfigProvider implementation of ConfigProvider
type defaultConfigProvider struct {
}

// ConfigProvider interface to implement to provide the app configuration
type ConfigProvider interface {
	GetApp() (*Config, error)
}

// DefaultSerializer returns the default App Serializer
func DefaultConfigProvider() ConfigProvider {
	return &defaultConfigProvider{}
}

// GetApp returns the app configuration
func (d *defaultConfigProvider) GetApp() (*Config, error) {

	configPath := config.GetFlogoConfigPath()

	flogo, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(flogo)
	app := &Config{}
	err = jsonParser.Decode(&app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func GetProperties(properties interface{}) (map[string]interface{}, error) {

	props := make(map[string]interface{})

	if properties != nil {
		oldPropsModel, ok := properties.(map[string]interface{})
		if ok {
			// Old model
			for name, value := range oldPropsModel {
				strValue, ok := value.(string)
				if ok {
					if strValue != "" && strValue[0] == '$' {
						// Needs resolution
						pValue, err := data.GetBasicResolver().Resolve(strValue, nil)
						if err != nil {
							return props, err
						}
						value = pValue
					}
				}
				props[name] = value
			}
			return props, nil
		}

		newPropModel, ok := properties.([]interface{})
		if ok {
			// New model
			for _, value := range newPropModel {
				propObj, ok := value.(map[string]interface{})
				if ok {
					pName := propObj["name"].(string)
					pType := propObj["type"].(string)
					pValue := propObj["value"]

					dataType, _ := data.ToTypeEnum(pType)
					strValue, ok := pValue.(string)
					if ok {
						if strValue != "" && strValue[0] == '$' {
							// Needs resolution
							resolvedValue, err := data.GetBasicResolver().Resolve(strValue, nil)
							if err != nil {
								return props, err
							}
							pValue = resolvedValue
						}
					}
					value, err := data.CoerceToValue(pValue, dataType)
					if err != nil {
						return props, err
					}
					props[pName] = value
				}
			}
			return props, nil
		}
		return nil, errors.New("Invalid application properties configuration")
	}

	return props, nil
}

func FixUpApp(cfg *Config) {

	if cfg.Resources != nil || cfg.Actions == nil {
		//already new app format
		return
	}

	idToAction := make(map[string]*action.Config)
	for _, act := range cfg.Actions {
		idToAction[act.Id] = act
	}

	for _, trg := range cfg.Triggers {
		for _, handler := range trg.Handlers {

			oldAction := idToAction[handler.ActionId]

			newAction := &action.Config{Ref: oldAction.Ref}

			if oldAction != nil {
				newAction.Mappings = oldAction.Mappings
			} else {
				if handler.ActionInputMappings != nil {
					newAction.Mappings = &data.IOMappings{}
					newAction.Mappings.Input = handler.ActionInputMappings
					newAction.Mappings.Output = handler.ActionOutputMappings
				}
			}

			newAction.Data = oldAction.Data
			newAction.Metadata = oldAction.Metadata

			handler.Action = newAction
		}
	}

	cfg.Actions = nil
}
