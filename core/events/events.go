package events

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"sync"
	"runtime/debug"
)

type EventListenerFunc func(*EventContext) error

var eventListeners = make(map[string][]EventListenerFunc)

// Buffered channel
var eventQueue = make(chan *EventContext, 100)
var publisherRoutineStarted = false
var shutdown = make(chan bool)

var lock = &sync.RWMutex{}

// Registers listener for flow events
func RegisterEventListener(eventType string, fel EventListenerFunc) error {
	lock.Lock()
	defer lock.Unlock()
	eventListeners[eventType] = append(eventListeners[eventType], fel)
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
		close(shutdown)
		publisherRoutineStarted = false
	}
}

//  EventContext is a wrapper over specific event context
type EventContext struct {
	// Type of event
	eventType string
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
			publishEvent(event)
			lock.RUnlock()
		case <-shutdown:
			return
		}
	}
}

func publishEvent(fe *EventContext) {

	for _, fel := range eventListeners[fe.eventType] {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("Registered event handler failed to process event due to error - '%v' ", r)
					logger.Errorf("StackTrace: %s", debug.Stack())
				}
			}()
			err := fel(fe)
			if err != nil {
				logger.Errorf("Registered event handler failed to process event due to error - '%s' ", err.Error())
			}
		}()
	}
}

//TODO channel to be passed to actions
func PublishEvent(eType string, event interface{}) {
	listeners := eventListeners[eType]
	if len(listeners) > 0 {
		evtContext := &EventContext{event: event, eventType: eType}
		// Put event on the queue
		eventQueue <- evtContext
	}
}
