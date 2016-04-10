package ppsremote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("processprovider")

// RemoteProcessProvider is an implementation of ProcessProvider service
// that can access processes via URI
type RemoteProcessProvider struct {
	//todo: switch to LRU cache
	processCache map[string]*process.Definition
}

// NewRemoteProcessProvider creates a RemoteProcessProvider
func NewRemoteProcessProvider() *RemoteProcessProvider {

	var service RemoteProcessProvider
	service.processCache = make(map[string]*process.Definition)

	return &service
}

// Start implements util.Managed.Start()
func (pps *RemoteProcessProvider) Start() {
	// no-op
}

// Stop implements util.Managed.Stop()
func (pps *RemoteProcessProvider) Stop() {
	// no-op
}

// Init implements services.ProcessProviderService.Init()
func (pps *RemoteProcessProvider) Init(settings map[string]string) {
	// no-op
}

// GetProcess implements process.Provider.GetProcess
func (pps *RemoteProcessProvider) GetProcess(processURI string) *process.Definition {

	if process, ok := pps.processCache[processURI]; ok {
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

	pps.processCache[processURI] = def

	return def
}
