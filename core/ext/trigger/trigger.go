package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
)

// Trigger is object that triggers/starts process instances and
// is managed by an engine
type Trigger interface {

	// TriggerMetadata returns the metadata of the trigger
	Metadata() *Metadata

	// Init sets up the trigger, it is called before Start()
	Init(starter processinst.Starter, config *Config)

	// Start starts the trigger
	Start()

	// Stop stops the trigger
	Stop()
}
