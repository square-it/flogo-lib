package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
)

// PooledRunner is a flow runner that queues and runs a flow in a worker pool
// todo: rename to AsyncFlowRunner?
type PooledRunner struct {
	workerQueue chan chan WorkRequest
	workQueue   chan WorkRequest
	numWorkers  int
	workers     []*Worker
	active      bool

	directRunner *DirectRunner
}

// PooledConfig is the configuration object for a PooledRunner
type PooledConfig struct {
	NumWorkers    int `json:"numWorkers"`
	WorkQueueSize int `json:"workQueueSize"`
	MaxStepCount  int `json:"maxStepCount"`
}

// NewPooledRunner create a new pooled
func NewPooledRunner(config *PooledConfig, stateRecorder flowinst.StateRecorder) *PooledRunner {

	var pooledRunner PooledRunner
	pooledRunner.directRunner = NewDirectRunner(stateRecorder, config.MaxStepCount)

	// config via engine config
	pooledRunner.numWorkers = config.NumWorkers
	pooledRunner.workQueue = make(chan WorkRequest, config.WorkQueueSize)

	return &pooledRunner
}

// Start will start the engine, by starting all of its workers
func (runner *PooledRunner) Start() error {

	if !runner.active {

		runner.workerQueue = make(chan chan WorkRequest, runner.numWorkers)

		runner.workers = make([]*Worker, runner.numWorkers)

		for i := 0; i < runner.numWorkers; i++ {
			log.Debug("Starting worker", i+1)
			worker := NewWorker(i+1, runner.directRunner, runner.workerQueue)
			runner.workers[i] = &worker
			worker.Start()
		}

		go func() {
			for {
				select {
				case work := <-runner.workQueue:
					log.Debug("Received work requeust")

					//todo fix, this creates unbounded go func blocked waiting for worker queue
					go func() {
						worker := <-runner.workerQueue

						log.Debug("Dispatching work request")
						worker <- work
					}()
				}
			}
		}()

		runner.active = true
	}

	return nil
}

// Stop will stop the engine, by stopping all of its workers
func (runner *PooledRunner) Stop() {

	if runner.active {

		runner.active = false

		for _, worker := range runner.workers {
			log.Debug("Stopping worker", worker.ID)
			worker.Stop()
		}
	}
}

// RunInstance runs the specified Flow Instance until it is complete
// or it no longer has any tasks to execute
func (runner *PooledRunner) RunInstance(instance *flowinst.Instance) bool {

	if runner.active {

		work := WorkRequest{ReqType: RtRun, Request: instance}

		runner.workQueue <- work
		log.Debugf("Run Flow '%s' queued", instance.ID())

		return false
	}

	//reject start

	return false
}
