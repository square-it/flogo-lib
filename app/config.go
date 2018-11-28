package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"io/ioutil"
	"regexp"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Config is the configuration for the App
type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Properties []*data.Attribute  `json:"properties"`
	Channels   []string           `json:"channels"`
	Triggers   []*trigger.Config  `json:"triggers"`
	Resources  []*resource.Config `json:"resources"`
	Actions    []*action.Config   `json:"actions"`
}

var appName, appVersion string

// Returns name of the application
func GetName() string {
	return appName
}

// Returns version of the application
func GetVersion() string {
	return appVersion
}

// Sets name of the application (useful when embedding flogo.json in embeddedapp.go)
func SetName(name string) {
	appName = name
}

// Sets version of the application (useful when embedding flogo.json in embeddedapp.go)
func SetVersion(version string) {
	appVersion = version
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
	return LoadConfig("")
}

func LoadConfig(flogoJson string) (*Config, error) {
	app := &Config{}
	if flogoJson == "" {
		configPath := config.GetFlogoConfigPath()

		flogo, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}

		file, err := ioutil.ReadAll(flogo)
		if err != nil {
			return nil, err
		}

		updated, err := preprocessConfig(file)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(updated, &app)
		if err != nil {
			return nil, err
		}

	} else {
		updated, err := preprocessConfig([]byte(flogoJson))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(updated, &app)
		if err != nil {
			return nil, err
		}
	}
	appName = app.Name
	appVersion = app.Version
	return app, nil
}

func preprocessConfig(appJson []byte) ([]byte, error) {

	// For now decode secret values
	re := regexp.MustCompile("SECRET:[^\\\\\"]*")
	for _, match := range re.FindAll(appJson, -1) {
		decodedValue, err := resolveSecretValue(string(match))
		if err != nil {
			return nil, err
		}
		appstring := strings.Replace(string(appJson), string(match), decodedValue, -1)
		appJson = []byte(appstring)
	}

	return appJson, nil
}

func resolveSecretValue(encrypted string) (string, error) {
	encodedValue := string(encrypted[7:])
	decodedValue, err := data.GetSecretValueHandler().DecodeValue(encodedValue)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}

func GetProperties(properties []*data.Attribute) (map[string]interface{}, error) {

	props := make(map[string]interface{})
	if properties != nil {
		overriddenProps, err := loadExternalProperties(properties)
		if err != nil {
			return props, err
		}
		for _, property := range properties {
			pValue := property.Value()
			if newValue, ok := overriddenProps[property.Name()]; ok {
				pValue = newValue
			}
			value, err := data.CoerceToValue(pValue, property.Type())
			if err != nil {
				return props, err
			}
			props[property.Name()] = value
		}
		return props, nil
	}

	return props, nil
}

func loadExternalProperties(properties []*data.Attribute) (map[string]interface{}, error) {
	logger.Debug("Resolving properties...")

	var resolvers []string
	if config.GetAppConfigEnvVars() {
		resolvers = append(resolvers, "env")
	}
	if config.GetAppConfigProfiles() != "" {
		resolvers = append(resolvers, "file")
	}
	if config.GetAppConfigExternal() != "" {
		resolvers = append(resolvers, config.GetAppConfigExternal())
	}

	props := make(map[string]interface{})

	for _, resolverType := range resolvers {
		resolver := GetPropertyValueResolver(resolverType)

		if resolver == nil {
			errMag := fmt.Sprintf("Unsupported resolver type '%s'. Resolver not registered.", resolverType)
			return nil, errors.New(errMag)
		}

		for _, prop := range properties {
			if props[prop.Name()] != nil { // Ignore property already resolved by another resolver
				continue
			}
			newVal, _ := resolver.ResolveValue(prop.Name())
			if newVal != nil { // if resolver returns nil, default value from flogo.json will be used
				logger.Debugf("Property '%s' resolved using resolver '%s'.", prop.Name(), resolverType)

				// Use new value
				props[prop.Name()] = newVal
				// May be a secret??
				strVal, _ := newVal.(string)
				if len(strVal) > 0 && strings.HasPrefix(strVal, "SECRET:") {
					// Resolve secret value
					newVal, err := resolveSecretValue(strVal)
					if err != nil {
						return nil, err
					}
					props[prop.Name()] = newVal
				}
			} else {
				logger.Debugf("Property '%s' could not be resolved using resolver '%s'.", prop.Name(), resolverType)
			}
		}
	}

	logger.Debug("Resolved properties")

	return props, nil
}

//used for old action config

//func FixUpApp(cfg *Config) {
//
//	if cfg.Resources != nil || cfg.Actions == nil {
//		//already new app format
//		return
//	}
//
//	idToAction := make(map[string]*action.Config)
//	for _, act := range cfg.Actions {
//		idToAction[act.Id] = act
//	}
//
//	for _, trg := range cfg.Triggers {
//		for _, handler := range trg.Handlers {
//
//			oldAction := idToAction[handler.ActionId]
//
//			newAction := &action.Config{Ref: oldAction.Ref}
//
//			if oldAction != nil {
//				newAction.Mappings = oldAction.Mappings
//			} else {
//				if handler.ActionInputMappings != nil {
//					newAction.Mappings = &data.IOMappings{}
//					newAction.Mappings.Input = handler.ActionInputMappings
//					newAction.Mappings.Output = handler.ActionOutputMappings
//				}
//			}
//
//			newAction.Data = oldAction.Data
//			newAction.Metadata = oldAction.Metadata
//
//			handler.Action = newAction
//		}
//	}
//
//	cfg.Actions = nil
//}
