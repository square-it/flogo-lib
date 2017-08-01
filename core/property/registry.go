package property

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"os"
)

var (
	props = make(map[string]interface{})
	mut   = sync.RWMutex{}
	regex = regexp.MustCompile(config.GetPropertyDelimiterFormat())
	// Default Resolver
	resolver Resolver = &DefaultResolver{}
)

// Resolve value sourced from Enviornment variable or any other configuration management services
type Resolver interface {
	//Resolve value for given name
	Resolve(name string) interface{}
}

type DefaultResolver struct {
	
}

// Default resolver resolves values from property bag and environment variable
func (resolver *DefaultResolver) Resolve(value string) interface{} {
	if len(value) == 0 {
		return value
	}
	// Value format: ${env.ENVVAR1}
	if strings.Contains(value, "${env.") {
		value = value[6 : len(value)-1]
		logger.Debugf("Resolving  value for enviornment variable: '%s'", value)
		return os.Getenv(value)
	} else if strings.Contains(value, "${property.") {
		// Value format: ${property.Prop1}
		// This is property bag resolution
		property := value[11 : len(value)-1]
		return Get(property)
	}
	return value
}

// Get returns the value of the property for the given id
func Get(id string) interface{} {
	mut.RLock()
	defer mut.RUnlock()
	prop, ok := props[id]
	if !ok {
		// Use resolver to resolve value
		return prop
	}

	switch prop.(type) {
	case string:
		value := prop.(string)
		// further resolution needed?
		if regex.MatchString(value) {
			if resolver != nil {
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

func RegisterResolver(newresolver Resolver) {
	mut.Lock()
	defer mut.Unlock()

	logger.Debugf("Registering property resolver: '%s'", reflect.TypeOf(newresolver).String())
	resolver = newresolver
}

func GetResolver() Resolver {
	return resolver
}
