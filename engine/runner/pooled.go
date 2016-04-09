package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
)

// PooledRunner is a process runner that queues and runs a process in a worker pool
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
func NewPooledRunner(config *PooledConfig, stateRecorder processinst.StateRecorder) *PooledRunner {

	var pooledRunner PooledRunner
	pooledRunner.directRunner = NewDirectRunner(stateRecorder, config.MaxStepCount)

	// config via engine config
	pooledRunner.numWorkers = config.NumWorkers
	pooledRunner.workQueue = make(chan WorkRequest, config.WorkQueueSize)

	return &pooledRunner
}

// Start will start the engine, by starting all of its workers
func (runner *PooledRunner) Start() {

	if !runner.active {

		log.Debug("Starting Pooled Process Instance Runner...")

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
					go func() {
						worker := <-runner.workerQueue

						log.Debug("Dispatching work request")
						worker <- work
					}()
				}
			}
		}()

		runner.active = true

		log.Debug("Started Pooled Process Instance Runner")
	}
}

// Stop will stop the engine, by stopping all of its workers
func (runner *PooledRunner) Stop() {

	if runner.active {

		log.Debug("Stopping Pooled Process Instance Runner...")

		runner.active = false

		for _, worker := range runner.workers {
			log.Debug("Stopping worker", worker.ID)
			worker.Stop()
		}

		log.Debug("Stopped Pooled Process Instance Runner")
	}
}

// RunInstance runs the specified Process Instance until it is complete
// or it no longer has any tasks to execute
func (runner *PooledRunner) RunInstance(instance *processinst.Instance) bool {

	if runner.active {

		work := WorkRequest{ReqType: RtRun, Request: instance}

		runner.workQueue <- work
		log.Debugf("Run Process '%s' queued", instance.ID())

		return false
	}

	//reject start

	return false
}
