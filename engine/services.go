package engine

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"bytes"
)

/////////////////////////////////////////////////////////////////////////////
// Process Registry - Implements process.Provider
// todo: make this pluggable

type ProcessRegistry struct {
	Processes map[string]*process.Definition
}

func NewProcessRegistry() *ProcessRegistry {

	var processRegistry ProcessRegistry
	processRegistry.Processes = make(map[string]*process.Definition)

	return &processRegistry
}

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

func (pr *ProcessRegistry) RegisterProcess(uri string, process *process.Definition) {

	log.Debugf("registering process: %s\n", uri)

	pr.Processes[uri] = process
}

/////////////////////////////////////////////////////////////////////////////
// REST StateService - Implements StateRecorder

type RestStateService struct {
	host string
}

func NewRestStateService(host string) processinst.StateRecorder {

	var stateService RestStateService
	stateService.host = host

	return &stateService
}

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

type RecordSnapshotReq struct {
	ID        int    `json:"id"`
	ProcessID string `json:"processID"`
	State     int    `json:"state"`
	Status    int    `json:"status"`

	SnapshotData *processinst.Instance `json:"snapshotData"`
}

type RecordStepReq struct {
	ID        int    `json:"id"`
	ProcessID string `json:"processID"`
	State     int    `json:"state"`
	Status    int    `json:"status"`

	StepData *processinst.InstanceChangeTracker `json:"stepData"`
}
