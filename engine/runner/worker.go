package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
)

// Based off: http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

// RequestType is value that indicates the type of Request
type RequestType int

const (
	// RtRun denotes a run flow instance request
	RtRun RequestType = 10
)

// WorkRequest describes a Request that Worker should handle
type WorkRequest struct {
	ReqType RequestType
	ID      string
	Request interface{}
}

// A Worker handles WorkRequest, work requests consist of start, restart
// and resume of FlowInstances
type Worker struct {
	ID          int
	runner      *DirectRunner
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, runner *DirectRunner, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		runner:      runner,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Start function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.  This is where all the request are handled
func (w Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				log.Debugf("worker-%d: Received Request\n", w.ID)

				switch work.ReqType {
				case RtRun:

					instance := work.Request.(*flowinst.Instance)

					log.Debugf("worker-%d: Received Run Request: %s - %s\n", w.ID, work.ID, instance.ID())
					w.runner.RunInstance(instance)
				}

			case <-w.QuitChan:
				// We have been asked to stop.
				log.Debugf("worker-%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
