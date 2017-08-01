package property

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var (
	props    = make(map[string]interface{})
	mut      = sync.RWMutex{}
	regex    = regexp.MustCompile(config.GetPropertyDelimiterFormat())
	resolver Resolver
)

func init() {
	logger.Debugf("Registering environment variable value resolver")
	RegisterResolver(&data.EnvVarResolver{})
}

// Resolve value sourced from Enviornment variable or any other configuration management services
type Resolver interface {
	//Resolve value for given name
	Resolve(name string) interface{}
}

// Get returns the value of the property for the given id
func Get(id string) interface{} {
	mut.RLock()
	defer mut.RUnlock()
	prop, ok := props[id]
	if !ok {
		return prop
	}

	switch prop.(type) {
	case string:
		value := prop.(string)
		// further resolution needed?
		if regex.MatchString(value) {
			if resolver != nil {
				//Value resolved by first resolver will be returned
				resolvedValue := resolver.Resolve(value)
				if resolvedValue != nil {
					logger.Debugf("Value is resolved by: '%s'", reflect.TypeOf(resolver).String())
					return resolvedValue
				}
			}
		}
		// Its literal value
		return value
	}
	return prop
}

func Register(id string, value interface{}) error {
	mut.Lock()
	defer mut.Unlock()

	if len(id) == 0 {
		return fmt.Errorf("error registering property, id is empty")
	}

	if _, ok := props[id]; ok {
		return fmt.Errorf("Error registering property, property already registered for id '%s'", id)
	}

	logger.Debugf("Registering property id: '%s', value: '%s'", id, value)

	props[id] = value

	return nil
}

// Use resolver to resolve values like ${MYVALUE}
func Resolve(value string) interface{} {
	mut.RLock()
	defer mut.RUnlock()
	// Should match with ${*}
	if regex.MatchString(value) {
		if strings.Contains(value, `{property.`) {
			// This is property bag resolution
			property := value[11 : len(value)-1]
			return Get(property)
		} else {
			// Call resolver to resolve value
			if resolver != nil {
				resolvedValue := resolver.Resolve(value)
				if resolvedValue != nil {
					logger.Debugf("Value is resolved by: '%s'", reflect.TypeOf(resolver).String())
					return resolvedValue
				}
			}
		}
	}
	// Return same value
	return value
}
func RegisterResolver(newresolver Resolver) {
	mut.Lock()
	defer mut.Unlock()

	logger.Debugf("Registering property resolver: '%s'", reflect.TypeOf(newresolver).String())
	resolver = newresolver
}
