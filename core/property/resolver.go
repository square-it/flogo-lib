package property

import (
	"os"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func init() {
	data.SetResolver(data.RES_PROPERTY, ResolveProperty)
	data.SetResolver(data.RES_ENV, ResolveEnv)
}

func ResolveProperty(scope data.Scope, value string) (interface{}, bool) {
	return Resolve(value)
}

func ResolveEnv(scope data.Scope, value string) (interface{}, bool) {
	if len(value) == 0 {
		return value, false
	}

	logger.Debugf("Resolving  value for environment variable: '%s'", value)
	return os.LookupEnv(value)
}

// Resolve interface, resolves a value for a given name
type Resolver interface {
	// Returns resolved value and true otherwise nil and false in case it can not resolve the value.
	Resolve(name string) (interface{}, bool)
}

// Default resolver to resolve property t variable values.
type DefaultPropertyResolver struct {
}

// Default resolver resolves values from property bag
func (resolver *DefaultPropertyResolver) Resolve(value string) (interface{}, bool) {
	if len(value) == 0 {
		return value, false
	}

	logger.Debugf("Resolving  value for property : '%s'", value)
	return Get(value)
}

