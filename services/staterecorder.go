package services

import (
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"encoding/json"
	"net/http"
	"bytes"
)

//todo create local,remote version of service

// RestStateService simple REST implementation of the StateRecorder
type RestStateService struct {
	host string
}

// NewRestStateService creates a new RestStateService
func NewRestStateService() StateRecorderService {

	return &RestStateService{}
}

func (ss *RestStateService) Start() {
	// no-op
}

func (ss *RestStateService) Stop() {
	// no-op
}

func (ss *RestStateService) Init(settings map[string]string) {
	ss.host = settings["host"]
}

// RecordSnapshot implements processinst.StateRecorder.RecordSnapshot
func (ss *RestStateService) RecordSnapshot(instance *processinst.Instance) {

	storeReq := &RecordSnapshotReq{
		ID:           instance.StepID(),
		ProcessID:    instance.ID(),
		State:        instance.State(),
		Status:       int(instance.Status()),
		SnapshotData: instance,
	}

	uri := ss.host + "/instances/snapshot"

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
func (ss *RestStateService) RecordStep(instance *processinst.Instance) {

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