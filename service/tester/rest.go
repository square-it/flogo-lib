package tester

import (
	"encoding/json"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/julienschmidt/httprouter"
)

// RestEngineTester is default REST implementation of the EngineTester
type RestEngineTester struct {
	reqProcessor *RequestProcessor
	server       *Server
	runner       runner.Runner
}

// NewRestEngineTester creates a new REST EngineTester
func NewRestEngineTester() *RestEngineTester {
	return &RestEngineTester{}
}

// Init implements engine.EngineTester.Init
func (et *RestEngineTester) Init(settings map[string]string) {

	router := httprouter.New()
	router.OPTIONS("/flow/start", handleOption)
	router.POST("/flow/start", et.StartFlow)

	router.OPTIONS("/flow/restart", handleOption)
	router.POST("/flow/restart", et.RestartFlow)

	router.OPTIONS("/flow/resume", handleOption)
	router.POST("/flow/resume", et.ResumeFlow)

	addr := ":" + settings["port"]
	et.server = NewServer(addr, router)
}

// Init implements engine.EngineTester.SetupInstanceSupport
func (et *RestEngineTester) SetupInstanceSupport(instManager *flowinst.Manager, runner runner.Runner) {

	et.reqProcessor = NewRequestProcessor(instManager)
	et.runner = runner
}

// Start implements engine.EngineTester.Start
func (et *RestEngineTester) Start() error {
	err := et.server.Start()
	return err
}

// Stop implements engine.EngineTester.Stop
func (et *RestEngineTester) Stop() {
	et.server.Stop()
}

func handleOption(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Origin")
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With")
	w.Header().Add("Access-Control-Allow-Headers", "Accept")
	w.Header().Add("Access-Control-Allow-Headers", "Accept-Language")
	w.Header().Set("Content-Type", "application/json")
}

// IDResponse is a respone object consists of an ID
type IDResponse struct {
	ID string `json:"id"`
}

// StartFlow starts a new Flow Instance (POST "/flow/start").
//
// To post a start flow, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{"flowUri":"base"}' http://localhost:8080/flow/start
func (et *RestEngineTester) StartFlow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	req := &StartRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	instance, err := et.reqProcessor.StartFlow(req, nil) //nil replyHandler

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &IDResponse{ID: instance.ID()}

	log.Debugf("Starting Instance [ID:%s] for %s", instance.ID(), req.FlowURI)

	et.runner.RunInstance(instance)

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// RestartFlow restarts a Flow Instance (POST "/flow/restart").
//
// To post a restart flow, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{...}' http://localhost:8080/flow/restart
func (et *RestEngineTester) RestartFlow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Error("Unable to restart flow, make sure definition registered")
	//	}
	//}()

	req := &RestartRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	instance := et.reqProcessor.RestartFlow(req, nil) //nil replyHandler

	// If we didn'et find it, 404
	//w.WriteHeader(http.StatusNotFound)

	resp := &IDResponse{ID: instance.ID()}

	et.runner.RunInstance(instance)

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// ResumeFlow resumes a Flow Instance (POST "/flow/resume").
//
// To post a resume flow, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{...}' http://localhost:8080/flow/resume
func (et *RestEngineTester) ResumeFlow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer func() {
		if r := recover(); r != nil {
			log.Error("Unable to resume flow, make sure definition registered")
		}
	}()

	req := &ResumeRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	instance := et.reqProcessor.ResumeFlow(req, nil) //nil replyHandler
	et.runner.RunInstance(instance)

	w.WriteHeader(http.StatusOK)
}
