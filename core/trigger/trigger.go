package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/util/managed"
)

// Factory is used to create new instances for a trigger
type Factory interface {
	New(config *Config) Trigger
}

// Trigger is object that triggers/starts flow instances and
// is managed by an engine
type Trigger interface {
	managed.Managed

	// Metadata returns the metadata of the trigger
	Metadata() *Metadata

	// Init sets up the trigger, it is called before Start()
	//DEPRECATED
	Init(actionRunner action.Runner)
}

// Init interface should be implemented by all Triggers, the Initialize method
// will eventually move up to Trigger to replace the deprecated "Init" method
type Init interface {

	// Init initialize the Activity for a particular configuration
	Initialize(ctx InitContext) error
}

type InitContext interface {

	// GetHandlers gets the handlers associated with the trigger
	GetHandlers() []*Handler
}