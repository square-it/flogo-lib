package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// Runner runs a flow instance
// todo: rename to FlowRunner?
type Runner interface {
	util.Managed

	// RunInstance run the specified flow instance
	RunInstance(instance *flowinst.Instance) bool
}
