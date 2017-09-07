package property

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var (
	props = make(map[string]interface{})
	mut   = sync.RWMutex{}

	// Default Property Resolver
	resolver Resolver = &DefaultPropertyResolver{}
)


// Get retrieves the value of the property named by the id.
// If the property is present in the registry the value
// (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be nil and the boolean will
// be false.
func Get(id string) (interface{}, bool) {
	mut.RLock()
	defer mut.RUnlock()
	prop, ok := props[id]
	return prop, ok
}

// Resolve resolves the value expressions like ${property.Prop1}
// using registered resolver. If it can handle such resolution, the resolved
// value is returned and the boolean is true. Otherwise the returned value will
// be nil and the boolean will be false.
func Resolve(value string) (interface{}, bool) {
	mut.RLock()
	defer mut.RUnlock()
	return resolver.Resolve(value)
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
