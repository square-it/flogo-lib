package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// Runner runs a process instance
// todo: rename to ProcessRunner?
type Runner interface {
	util.Managed

	// RunInstance run the specified process instance
	RunInstance(instance *processinst.Instance) bool
}
