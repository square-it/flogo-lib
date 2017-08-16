package expr

// Resolve value sourced from Environment variable or any other configuration management services
type Resolver interface {
	// Resolve value for given name
	// Returns resolved value and true otherwise nil and false in case it can not resolve the value.
	Resolve(name string) (interface{}, bool)
}
