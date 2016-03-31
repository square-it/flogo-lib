package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/engine/starter"
)

// Trigger is object that triggers/starts process instances and
// is managed by an engine
type Trigger interface {

	// TriggerMetadata returns the metadata of the trigger
	Metadata() *Metadata

	// Init sets up the trigger, it is called before Start()
	// todo: switch to config map[string]interface{}
	Init(processStarter starter.ProcessStarter, config map[string]string)

	// Start starts the trigger
	Start()

	// Stop stops the trigger
	Stop()
}
