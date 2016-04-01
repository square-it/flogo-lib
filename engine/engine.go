package engine

import (
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/TIBCOSoftware/flogo-lib/engine/starter"
)


// Engine creates and executes ProcessInstances.  It contains
// a Worker pool that handles the instance executions.
type Engine struct {
	workerQueue chan chan WorkRequest
	workQueue   chan WorkRequest
	numWorkers  int
	workers     []*Worker
	active      bool
	runner      *Runner
	generator   *util.Generator

	//lock sync.Mutex
}

// NewEngine create a new Engine
func NewEngine(system *System, config *EngineConfig) *Engine {
	var engine Engine
	engine.numWorkers = config.NumWorkers
	engine.active = false
	engine.runner = NewRunner(system, config.MaxStepCount)
	engine.generator, _ = util.NewGenerator()

	// config via engine config
	engine.workQueue = make(chan WorkRequest, config.WorkQueueSize)

	return &engine
}

// Start will start the engine, by starting all of its workers
func (e *Engine) Start() {
	// e.lock.Lock()
	// defer e.lock.Unlock()

	if !e.active {

		e.workerQueue = make(chan chan WorkRequest, e.numWorkers)

		e.workers = make([]*Worker, e.numWorkers)

		for i := 0; i < e.numWorkers; i++ {
			log.Debug("Starting worker", i+1)
			worker := NewWorker(i+1, e.runner, e.workerQueue)
			e.workers[i] = &worker
			worker.Start()
		}

		go func() {
			for {
				select {
				case work := <-e.workQueue:
					log.Debug("Received work requeust")
					go func() {
						worker := <-e.workerQueue

						log.Debug("Dispatching work request")
						worker <- work
					}()
				}
			}
		}()

		e.active = true
	}
}

// Stop will stop the engine, by stopping all of its workers
func (e *Engine) Stop() {

	// e.lock.Lock()
	// defer e.lock.Unlock()
	if e.active {

		for _, worker := range e.workers {
			worker.Stop()
		}

		//stop loop
	}
}

// StartProcess implements engine.ProcessStarter.StartProcess
func (e *Engine) StartProcess(startRequest *starter.StartRequest) string {

	if e.active {

		instanceID := e.generator.NextAsString()

		//todo should we create instance here? so we can immediately return an error if necessary
		work := WorkRequest{ReqType: RtStart, ID: instanceID, Request: startRequest}

		e.workQueue <- work
		log.Debug("Start Process queued")

		return instanceID
	}

	//reject start

	return ""
}

// RestartProcess implements engine.ProcessStarter.RestartProcess
func (e *Engine) RestartProcess(restartRequest *starter.RestartRequest) string {

	if e.active {

		instanceID := e.generator.NextAsString()

		//todo should we create instance here? so we can immediately return an error if necessary
		work := WorkRequest{ReqType: RtRestart, ID: instanceID, Request: restartRequest}

		e.workQueue <- work
		log.Debug("Restart Process queued")

		return instanceID
	}

	//reject restart

	return ""
}

// ResumeProcess implements engine.ProcessStarter.ResumeProcess
func (e *Engine) ResumeProcess(resumeRequest *starter.ResumeRequest) {

	if e.active {

		work := WorkRequest{ReqType: RtResume, Request: resumeRequest}

		e.workQueue <- work
		log.Debug("Resume Process queued")
	} else {
		//reject resume
	}
}
