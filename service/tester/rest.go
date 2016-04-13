package tester

import (
	"encoding/json"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
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
func (et *RestEngineTester) Init(instManager *processinst.Manager, runner runner.Runner, config map[string]string) {

	et.reqProcessor = NewRequestProcessor(instManager)
	et.runner = runner

	router := httprouter.New()
	router.OPTIONS("/process/start", handleOption)
	router.POST("/process/start", et.StartProcess)

	router.OPTIONS("/process/restart", handleOption)
	router.POST("/process/restart", et.RestartProcess)

	router.OPTIONS("/process/resume", handleOption)
	router.POST("/process/resume", et.ResumeProcess)

	addr := ":" + config["port"]
	et.server = NewServer(addr, router)
}

// Start implements engine.EngineTester.Start
func (et *RestEngineTester) Start() {
	et.server.Start()
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

// StartProcess starts a new Process Instance (POST "/process/start").
//
// To post a start process, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{"processUri":"base"}' http://localhost:8080/process/start
func (et *RestEngineTester) StartProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	req := &StartRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	instance := et.reqProcessor.StartProcess(req, nil) //nil replyHandler

	// If we didn't find it, 404
	//w.WriteHeader(http.StatusNotFound)

	resp := &IDResponse{ID: instance.ID()}

	log.Debugf("Starting Instance [ID:%s] for %s", instance.ID(), req.ProcessURI)

	et.runner.RunInstance(instance)

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// RestartProcess restarts a Process Instance (POST "/process/restart").
//
// To post a restart process, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{...}' http://localhost:8080/process/restart
func (et *RestEngineTester) RestartProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Error("Unable to restart process, make sure definition registered")
	//	}
	//}()

	req := &RestartRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	instance := et.reqProcessor.RestartProcess(req, nil) //nil replyHandler

	// If we didn'et find it, 404
	//w.WriteHeader(http.StatusNotFound)

	resp := &IDResponse{ID: instance.ID()}

	et.runner.RunInstance(instance)

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// ResumeProcess resumes a Process Instance (POST "/process/resume").
//
// To post a resume process, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{...}' http://localhost:8080/process/resume
func (et *RestEngineTester) ResumeProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer func() {
		if r := recover(); r != nil {
			log.Error("Unable to resume process, make sure definition registered")
		}
	}()

	req := &ResumeRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	instance := et.reqProcessor.ResumeProcess(req, nil) //nil replyHandler
	et.runner.RunInstance(instance)

	w.WriteHeader(http.StatusOK)
}
