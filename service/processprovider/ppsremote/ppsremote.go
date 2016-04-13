package ppsremote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("processprovider")

// RemoteProcessProvider is an implementation of ProcessProvider service
// that can access processes via URI
type RemoteProcessProvider struct {
	//todo: switch to LRU cache
	processCache map[string]*process.Definition
	embeddedMgr  *util.EmbeddedFlowManager
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
func (pps *RemoteProcessProvider) Init(settings map[string]string, embeddedFlowMgr *util.EmbeddedFlowManager) {
	pps.embeddedMgr = embeddedFlowMgr
}

// GetProcess implements process.Provider.GetProcess
func (pps *RemoteProcessProvider) GetProcess(processURI string) *process.Definition {

	// todo turn pps.processCache to real cache
	if process, ok := pps.processCache[processURI]; ok {
		log.Debugf("Accessing cached Process: %s\n")
		return process
	}

	log.Debugf("Get Process: %s\n", processURI)

	var flowJSON []byte

	if strings.Index(processURI, "local://") == 0 {

		log.Debugf("Loading Embedded Flow: %s\n", processURI)
		flowJSON = pps.embeddedMgr.GetEmbeddedFlowJSON(processURI[8:])

	} else {

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
		flowJSON = body
	}

	if flowJSON != nil {
		var defRep process.DefinitionRep
		json.Unmarshal(flowJSON, &defRep)

		def := process.NewDefinition(&defRep)

		pps.processCache[processURI] = def

		return def
	}

	log.Debugf("Flow Not Found: %s\n", processURI)

	return nil
}
