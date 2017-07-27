package property

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"reflect"
	"sync"
)

var (
	props     = make(map[string]string)
	mut       = sync.RWMutex{}
	resolvers []PropertyValueResolver
)

// Resolve value sourced from Enviornment variable or any other configuration management services
type PropertyValueResolver interface {
	//Resolve value for given name
	Resolve(name string) string
}

func RegisterDefaultResolver() {
	RegisterValueResolver(&data.EnvVarResolver{})
}

// Get returns the value of the property for the given id
// If it is an environment property (for example {MY_PROP})
// The value will be looked up in the os environment
func Get(id string) string {
	mut.RLock()
	defer mut.RUnlock()
	prop, ok := props[id]
	if !ok {
		return prop
	}
	for _, resolver := range resolvers {
		//Value resolved by first resolver will be returned
		prop = resolver.Resolve(prop)
		if prop != "" {
			return prop
		}
	}
	return prop
}

func Register(id, value string) error {
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

func RegisterValueResolver(resolver PropertyValueResolver) {
	mut.Lock()
	defer mut.Unlock()

	logger.Debugf("Registering property resolver: '%s'", reflect.TypeOf(resolver).String())
	resolvers = append(resolvers, resolver)
}
