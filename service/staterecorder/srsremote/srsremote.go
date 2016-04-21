package srsremote

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("staterecorder")

// RemoteStateRecorder is an implementation of StateRecorder service
// that can access flowes via URI
type RemoteStateRecorder struct {
	host string
}

// NewRemoteStateRecorder creates a new RemoteStateRecorder
func NewRemoteStateRecorder() *RemoteStateRecorder {

	return &RemoteStateRecorder{}
}

// Start implements util.Managed.Start()
func (srs *RemoteStateRecorder) Start() error {
	// no-op
	return nil
}

// Stop implements util.Managed.Stop()
func (srs *RemoteStateRecorder) Stop() {
	// no-op
}

// Init implements services.StateRecorderService.Init()
func (srs *RemoteStateRecorder) Init(settings map[string]string) {

	host, set := settings["host"]
	port, set := settings["port"]

	if !set {
		panic("RemoteStateRecorder: requried setting 'host' not set")
	}

	if strings.Index(host, "http") != 0 {
		srs.host = "http://" + host + ":" + port
	} else {
		srs.host = host + ":" + port
	}

	log.Debugf("RemoteStateRecorder: StateRecoder Server = %s", srs.host)
}

// RecordSnapshot implements flowinst.StateRecorder.RecordSnapshot
func (srs *RemoteStateRecorder) RecordSnapshot(instance *flowinst.Instance) {

	storeReq := &RecordSnapshotReq{
		ID:           instance.StepID(),
		FlowID:       instance.ID(),
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

// RecordStep implements flowinst.StateRecorder.RecordStep
func (ss *RemoteStateRecorder) RecordStep(instance *flowinst.Instance) {

	storeReq := &RecordStepReq{
		ID:       instance.StepID(),
		FlowID:   instance.ID(),
		State:    instance.State(),
		Status:   int(instance.Status()),
		StepData: instance.ChangeTracker,
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
	ID     int    `json:"id"`
	FlowID string `json:"flowID"`
	State  int    `json:"state"`
	Status int    `json:"status"`

	SnapshotData *flowinst.Instance `json:"snapshotData"`
}

// RecordStepReq serializable representation of the RecordStep request
type RecordStepReq struct {
	ID     int    `json:"id"`
	FlowID string `json:"flowID"`
	State  int    `json:"state"`
	Status int    `json:"status"`

	StepData *flowinst.InstanceChangeTracker `json:"stepData"`
}
