package srsremote

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("staterecorder")

// RemoteStateRecorder is an implementation of StateRecorder service
// that can access processes via URI
type RemoteStateRecorder struct {
	host string
}

// NewRemoteStateRecorder creates a new RemoteStateRecorder
func NewRemoteStateRecorder() *RemoteStateRecorder {

	return &RemoteStateRecorder{}
}

// Start implements util.Managed.Start()
func (srs *RemoteStateRecorder) Start() {
	// no-op
}

// Stop implements util.Managed.Stop()
func (srs *RemoteStateRecorder) Stop() {
	// no-op
}

// Init implements services.StateRecorderService.Init()
func (srs *RemoteStateRecorder) Init(settings map[string]string) {

	host, set := settings["host"]

	if !set {
		panic("RemoteStateRecorder: requried setting 'host' not set")
	}

	srs.host = host
}

// RecordSnapshot implements processinst.StateRecorder.RecordSnapshot
func (srs *RemoteStateRecorder) RecordSnapshot(instance *processinst.Instance) {

	storeReq := &RecordSnapshotReq{
		ID:           instance.StepID(),
		ProcessID:    instance.ID(),
		State:        instance.State(),
		Status:       int(instance.Status()),
		SnapshotData: instance,
	}

	uri := srs.host + "/instances/snapshot"

	log.Debugf("POST Snapshot: %s\n", uri)

	jsonReq, _ := json.Marshal(storeReq)

	log.Debug("JSON: ", string(jsonReq))

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		//error
	}
}

// RecordStep implements processinst.StateRecorder.RecordStep
func (ss *RemoteStateRecorder) RecordStep(instance *processinst.Instance) {

	storeReq := &RecordStepReq{
		ID:        instance.StepID(),
		ProcessID: instance.ID(),
		State:     instance.State(),
		Status:    int(instance.Status()),
		StepData:  instance.ChangeTracker,
	}

	uri := ss.host + "/instances/steps"

	log.Debugf("POST Snapshot: %s\n", uri)

	jsonReq, _ := json.Marshal(storeReq)

	log.Debug("JSON: ", string(jsonReq))

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		//error
	}
}

// RecordSnapshotReq serializable representation of the RecordSnapshot request
type RecordSnapshotReq struct {
	ID        int    `json:"id"`
	ProcessID string `json:"processID"`
	State     int    `json:"state"`
	Status    int    `json:"status"`

	SnapshotData *processinst.Instance `json:"snapshotData"`
}

// RecordStepReq serializable representation of the RecordStep request
type RecordStepReq struct {
	ID        int    `json:"id"`
	ProcessID string `json:"processID"`
	State     int    `json:"state"`
	Status    int    `json:"status"`

	StepData *processinst.InstanceChangeTracker `json:"stepData"`
}
