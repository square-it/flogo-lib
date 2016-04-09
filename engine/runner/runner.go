package runner

import "github.com/TIBCOSoftware/flogo-lib/core/processinst"

// Runner runs a process instance
type Runner interface {
	// Start starts the runner
	Start()

	// Stop stops the runner
	Stop()

	// RunInstance run the specified process instance
	RunInstance(instance *processinst.Instance) bool
}
