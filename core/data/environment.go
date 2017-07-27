package data

import (
	"os"
)

type EnvVarResolver struct {
}

// This function will check if the value is an environment value (for example {MY_VALUE})
// if it is an environment value it will get resolved, otherwise the original value is returned
func (resolver *EnvVarResolver) Resolve(value string) interface{} {
	if len(value) == 0 {
		return value
	}
	return os.Getenv(value)
}
