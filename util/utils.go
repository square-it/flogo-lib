package util

import (
	"fmt"
	"runtime/debug"

	"github.com/op/go-logging"
)

// HandlePanic helper method to handl panics
func HandlePanic(name string, err *error) {
	if r := recover(); r != nil {

		log.Warningf("%s: PANIC Occurred  : %v\n", name, r)

		// todo: useful for debugging
		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("StackTrace: %s", debug.Stack())
		}

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
