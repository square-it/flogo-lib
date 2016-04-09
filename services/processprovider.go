package services

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

//todo create local,remote version of service

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

func (pr *ProcessRegistry) Start() {
	// no-op
}

func (pr *ProcessRegistry) Stop() {
	// no-op
}

func (pr *ProcessRegistry) Init(settings map[string]string) {
	// no-op
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
