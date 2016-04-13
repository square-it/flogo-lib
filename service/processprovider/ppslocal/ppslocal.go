package ppslocal

import (
	"encoding/json"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("processprovider")

// LocalProcessProvider is an implementation of ProcessProvider service
// that can access processes via URI
type LocalProcessProvider struct {
	//todo: switch to LRU cache
	processCache map[string]*process.Definition
	embeddedMgr  *util.EmbeddedFlowManager
}

// NewLocalProcessProvider creates a LocalProcessProvider
func NewLocalProcessProvider() *LocalProcessProvider {

	var service LocalProcessProvider
	service.processCache = make(map[string]*process.Definition)

	return &service
}

// Start implements util.Managed.Start()
func (pps *LocalProcessProvider) Start() {
	// no-op
}

// Stop implements util.Managed.Stop()
func (pps *LocalProcessProvider) Stop() {
	// no-op
}

// Init implements services.ProcessProviderService.Init()
func (pps *LocalProcessProvider) Init(settings map[string]string, embeddedFlowMgr *util.EmbeddedFlowManager) {
	pps.embeddedMgr = embeddedFlowMgr
}

// GetProcess implements process.Provider.GetProcess
func (pps *LocalProcessProvider) GetProcess(processURI string) *process.Definition {

	// todo turn pps.processCache to real cache
	if process, ok := pps.processCache[processURI]; ok {
		log.Debugf("Accessing cached Process: %s\n")
		return process
	}

	log.Debugf("Get Flow: %s\n", processURI)

	var flowJSON []byte

	if strings.Index(processURI, "local://") == 0 {

		flowJSON = pps.embeddedMgr.GetEmbeddedFlowJSON(processURI[8:])

	} else {

		flowJSON = pps.embeddedMgr.GetEmbeddedFlowJSON(processURI)
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
