package data

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"os"
	"strings"
)

type EnvVarResolver struct {
}


// This function will check if the value is an environment value (for example ${env.MY_VALUE})
// if it is an environment value it will get resolved, otherwise the original value is returned
func (resolver *EnvVarResolver) Resolve(value string) interface{} {
	if len(value) == 0 {
		return value
	}
	if strings.Contains(value, "${env.") {
		value = value[6 : len(value)-1]
		logger.Debugf("Resolving  value for enviornment variable: '%s'", value)
		return os.Getenv(value)
	}
	return value
}
