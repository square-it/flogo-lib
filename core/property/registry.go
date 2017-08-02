package property

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"os"
	"reflect"
	"strings"
	"sync"
)

var (
	props = make(map[string]interface{})
	mut   = sync.RWMutex{}
	// Default Resolver
	resolver Resolver = &DefaultResolver{}
)

// Resolve value sourced from Enviornment variable or any other configuration management services
type Resolver interface {
	// Resolve value for given name
	// Returns resolved value and true ortherwise nil and false in case it can not resolve the value.
	Resolve(name string) (interface{}, bool)
}

// Default resolver to resolve property and envioronment variable values.
type DefaultResolver struct {
}

// Default resolver resolves values from property bag and environment variable
func (resolver *DefaultResolver) Resolve(value string) (interface{}, bool) {
	if len(value) == 0 {
		return value, true
	}

	if strings.Contains(value, "${env.") {
		// Value format: ${env.ENVVAR1}
		key := value[6 : len(value)-1]
		logger.Debugf("Resolving  value for enviornment variable: '%s'", key)
		return os.LookupEnv(key)
	} else if strings.Contains(value, "${property.") {
		// Value format: ${property.Prop1}
		property := value[11 : len(value)-1]
		logger.Debugf("Resolving  value for property : '%s'", property)
		return Get(property)
	}

	// No resolution needed
	return value, true
}

// Get retrieves the value of the proprty named by the id.
// If the property is present in the registry the value
// (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be nil and the boolean will
// be false.
func Get(id string) (interface{}, bool) {
	mut.RLock()
	defer mut.RUnlock()
	prop, ok := props[id]
	if !ok {
		return prop, ok
	}
	return getValueFromResolver(prop)
}

func getValueFromResolver(prop interface{}) (interface{}, bool) {
	switch prop.(type) {
	case string:
		value := prop.(string)
		// Resolver can resolve the value
		resolvedValue, ok := resolver.Resolve(value)
		if ok {
			logger.Debugf("Value is resolved by: '%s'", reflect.TypeOf(resolver).String())
			return resolvedValue, ok
		}
		return prop, false
	}
	return prop, true
}

// Resolve resolves the value expressions like ${property.Prop1}
// or ${env.ENVVAR} using registered resolver. If it can handle such
// expressions, the resolved value is returned and the boolean is true.
// Otherwise the returned value will be nil and the boolean will be false.

func Resolve(name string) (interface{}, bool) {
	mut.RLock()
	defer mut.RUnlock()
	return resolver.Resolve(name)
}

// Register property with given value.
// Only String, Boolean and integer values are supported.
func Register(id string, value interface{}) error {
	mut.Lock()
	defer mut.Unlock()

	if len(id) == 0 {
		return fmt.Errorf("error registering property, id is empty")
	}

	if _, ok := props[id]; ok {
		return fmt.Errorf("Error registering property, property already registered for id '%s'", id)
	}

	switch value.(type) {
	case string:
	case bool:
	case int64:
	default:
		return fmt.Errorf("Error registering property: '%s'. Unsupported type: '%T'. Supported Types: [string, bool, int64]", id, value)
	}

	logger.Debugf("Registering property id: '%s', value: '%s'", id, value)

	props[id] = value

	return nil
}

// Register new resolver with the engine.
// It will override default resolver.
func RegisterResolver(newresolver Resolver) {
	if newresolver != nil {
		mut.Lock()
		defer mut.Unlock()
		logger.Debugf("Registering property resolver: '%s'", reflect.TypeOf(newresolver).String())
		resolver = newresolver
	}
}
