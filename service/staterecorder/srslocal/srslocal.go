package srslocal

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("staterecorder")

// LocalStateRecorder is an implementation of StateRecorder service
// that stores the state locally
type LocalStateRecorder struct {
	host string
}

// NewLocalStateRecorder creates a new LocalStateRecorder
func NewLocalStateRecorder() *LocalStateRecorder {

	return &LocalStateRecorder{}
}

// Start implements util.Managed.Start()
func (srs *LocalStateRecorder) Start() {
	// no-op
}

// Stop implements util.Managed.Stop()
func (srs *LocalStateRecorder) Stop() {
	// no-op
}

// Init implements services.StateRecorderService.Init()
func (srs *LocalStateRecorder) Init(settings map[string]string) {
	//srs.host = settings["host"]
}

// RecordSnapshot implements flowinst.StateRecorder.RecordSnapshot
func (srs *LocalStateRecorder) RecordSnapshot(instance *flowinst.Instance) {

}

// RecordStep implements flowinst.StateRecorder.RecordStep
func (srs *LocalStateRecorder) RecordStep(instance *flowinst.Instance) {

}
