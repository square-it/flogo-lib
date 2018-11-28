package app

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"os"
	"strings"
)

func init() {
	RegisterPropertyValueResolver("env", &EnvVariableValueResolver{})
}

// Resolve property value from environment variable
type EnvVariableValueResolver struct {

}

func (resolver *EnvVariableValueResolver) ResolveValue(key string) (interface{}, error) {
	value, exists := os.LookupEnv(key) // first try with the name of the property as is
	if exists {
		return value, nil
	}

	value, exists = os.LookupEnv(getCanonicalEnv(key)) // if not found try with the canonical form

	if !exists {
		return nil, errors.New(fmt.Sprintf("Environment variable '%s' is not set", key))
	}
	return value, nil
}

func getCanonicalEnv(key string) string {
	result := strcase.ToScreamingSnake(key)
	result = strings.Replace(result, ".", "_", -1)
	result = strings.Replace(result, "__", "_", -1)
	return result
}
