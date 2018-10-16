package events

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"sync"
	"errors"
	"runtime/debug"
)

type EventListenerFunc func(*EventContext) error

var eventListeners = make(map[string]EventListenerFunc)

// Buffered channel
var eventQueue = make(chan interface{}, 100)
var publisherRoutineStarted = false
var shutdown = make(chan bool)

var lock = &sync.RWMutex{}

// Registers listener for flow events
func RegisterEventListener(name string, fel EventListenerFunc) error {
	lock.Lock()
	defer lock.Unlock()
	_, exists := eventListeners[name]
	if exists {
		errMsg := fmt.Sprintf("Event listener with name - '%s' already registered", name)
		logger.Error(errMsg)
		return errors.New(errMsg)
	}
	eventListeners[name] = fel
	startPublisherRoutine()
	return nil
}

// Unregisters flow event listener
func UnRegisterEventListener(name string) {
	lock.Lock()
	defer lock.Unlock()
	_, exists := eventListeners[name]
	if exists {
		delete(eventListeners, name)
	}
	stopPublisherRoutine()
}

func startPublisherRoutine() {
	if publisherRoutineStarted == true {
		return
	}

	if len(eventListeners) > 0 {
		// start publisher routine
		go publishEvents()
		publisherRoutineStarted = true
	}
}

func stopPublisherRoutine() {
	if publisherRoutineStarted == false {
		return
	}

	if len(eventListeners) == 0 {
		// No more listeners. Stop go routine
		shutdown <- true
		publisherRoutineStarted = false
	}
}

//  EventContext is a wrapper over specific event context
type EventContext struct {
	// Event can be FlowEventContext or TaskEventContext
	event interface{}
}

// Returns wrapped event
func (ec *EventContext) GetEvent() interface{} {
	return ec.event
}

func publishEvents() {
	for {
		select {
		case event := <-eventQueue:
			lock.RLock()
			evtContext := &EventContext{event: event}
			publishEvent(evtContext)
			lock.RUnlock()
		case <-shutdown:
			return
		}
	}
}

func publishEvent(fe *EventContext) {

	for name, fel := range eventListeners {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("Registered event handler - '%s' failed to process event due to error - '%v' ", name, r)
					logger.Errorf("StackTrace: %s", debug.Stack())
				}
			}()
			err := fel(fe)
			if err != nil {
				logger.Errorf("Registered event handler - '%s' failed to process event due to error - '%s' ", name, err.Error())
			}
		}()
	}
}

//TODO channel to be passed to actions
// Put event on the queue
func PublishEvent(event interface{}) {
	eventQueue <- event
}