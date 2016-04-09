package engine

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
)

/////////////////////////////////////////////////////////////////////////////
// Process Registry - Implements process.Provider
// todo: make this pluggable

// ProcessRegistry is a simple Process registry
type ProcessRegistry struct {
	Processes map[string]*process.Definition
}

// NewProcessRegistry creates a new Process registry
func NewProcessRegistry() *ProcessRegistry {

	var processRegistry ProcessRegistry
	processRegistry.Processes = make(map[string]*process.Definition)

	return &processRegistry
}

// GetProcess implements process.Provider.GetProcess
func (pr *ProcessRegistry) GetProcess(processURI string) *process.Definition {

	if process, ok := pr.Processes[processURI]; ok {
		log.Debugf("Accessing cached Process: %s\n")
		return process
	}

	log.Debugf("GET Process: %s\n", processURI)

	req, err := http.NewRequest("GET", processURI, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		//not found
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)

	result := string(body)

	var defRep process.DefinitionRep
	json.Unmarshal([]byte(result), &defRep)

	def := process.NewDefinition(&defRep)

	pr.Processes[processURI] = def

	return def
}

// RegisterProcess registers a process locally
func (pr *ProcessRegistry) RegisterProcess(uri string, process *process.Definition) {

	log.Debugf("registering process: %s\n", uri)

	pr.Processes[uri] = process
}

/////////////////////////////////////////////////////////////////////////////
// REST StateService - Implements StateRecorder

// RestStateService simple REST implementation of the StateRecorder
type RestStateService struct {
	host string
}

// NewRestStateService creates a new RestStateService
func NewRestStateService(host string) processinst.StateRecorder {

	var stateService RestStateService
	stateService.host = host

	return &stateService
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

/////////////////////////////////////////////////////////////////////////////
// Tester
// todo create Managed interface (with Start/Stop)

// Tester is an engine interface to assist in testing processes
type Tester interface {

	//Init initializes the EngineTester
	Init(instManager *processinst.Manager, runner runner.Runner, config map[string]string)

	// Start starts the EngineTester
	Start()

	// Stop stops the EngineTester
	Stop()
}
