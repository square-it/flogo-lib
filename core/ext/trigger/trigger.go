package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// Trigger is object that triggers/starts process instances and
// is managed by an engine
type Trigger interface {
	util.Managed

	// TriggerMetadata returns the metadata of the trigger
	Metadata() *Metadata

	// Init sets up the trigger, it is called before Start()
	Init(starter processinst.Starter, config *Config)
}
