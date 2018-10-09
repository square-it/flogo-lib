package app

import (
	"sync"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"fmt"
	"errors"
	"os"
)

var (
	propertyProvider   *PropertyProvider
	propValueResolvers = make(map[string]PropertyValueResolver)
	lock               = &sync.Mutex{}
)


func init() {
	propertyProvider = &PropertyProvider{properties: make(map[string]interface{})}
	RegisterPropertyValueResolver("env", &EnvVariableValueResolver{})
}

func RegisterPropertyValueResolver(relType string, resolver PropertyValueResolver) error {
	lock.Lock()
	defer lock.Unlock()
	_, ok := propValueResolvers[relType]
	if ok {
		errMsg := fmt.Sprintf("Property value resolver is already registered for type - '%s'", relType)
		logger.Errorf(errMsg)
		return errors.New(errMsg)
	}
	propValueResolvers[relType] = resolver
	return nil
}

func GetPropertyValueResolver(relType string) PropertyValueResolver {
    return propValueResolvers[relType]
}

// Resolve property value from environment variable
type EnvVariableValueResolver struct {
	
}

func (resolver *EnvVariableValueResolver) ResolveValue(toResolve string) (interface{}, error) {
	value, exists := os.LookupEnv(toResolve)
	if !exists {
		return nil, errors.New(fmt.Sprintf("Environment variable - %s is not set",toResolve))
	}
	return value, nil
}


func GetPropertyProvider() *PropertyProvider {
	return propertyProvider
}

type PropertyProvider struct {
	properties map[string]interface{}
}

type PropertyValueResolver interface {
	ResolveValue(toResolve string) (interface{}, error)
}

func (pp *PropertyProvider) GetProperty(property string) (interface{}, bool) {
	prop, exists := pp.properties[property]
	return prop, exists
}

func (pp *PropertyProvider) SetProperty(property string, value interface{}) {
	pp.properties[property] = value
}

func (pp *PropertyProvider) SetProperties(value map[string]interface{}) {
	pp.properties = value
}
