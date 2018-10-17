package events

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"sync"
	"runtime/debug"
	"errors"
	"strings"
)

type EventListener interface {
	// Returns name of the listener
	Name() string

	// Returns list of event types interested in
	EventTypes() []string

	// Called when matching event occurs
	HandleEvent(*EventContext) error
}

var eventListeners = make(map[string][]EventListener)

// Buffered channel
var eventQueue = make(chan *EventContext, 100)
var publisherRoutineStarted = false
var shutdown = make(chan bool)

var lock = &sync.RWMutex{}

// Registers listener for events
func RegisterEventListener(evtListener EventListener) error {
	if evtListener == nil {
		return errors.New("Event handler must not nil")
	}

	if len(evtListener.EventTypes()) == 0 {
		return errors.New("Failed register event handler. At-least one event type must be configured.")
	}

	lock.Lock()
	defer lock.Unlock()

	for _, eType := range evtListener.EventTypes() {
		eventListeners[eType] = append(eventListeners[eType], evtListener)
		logger.Debugf("Event Listener - '%s' successfully registered for event type - '%s'", evtListener.Name(), eType)
	}

	startPublisherRoutine()
	return nil
}

// Unregisters event listener for given name.
// Set eventType to unregister listener from specific event types
func UnRegisterEventListener(name string, eventTypes ...string) {

	if name == "" {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	var deleteList []string
	var index = -1

	if len(eventTypes) > 0 {
		for _, eType := range eventTypes {
			evtLs, ok := eventListeners[eType]
			if ok {
				for i, el := range evtLs {
					if strings.EqualFold(el.Name(), name) {
						index = i
						break
					}
				}
				if index > -1 {
					if len(evtLs) > 1 {
						// More than one listeners. Just adjust slice
						eventListeners[eType] = append(eventListeners[eType][:index], eventListeners[eType][index+1:]...)
					} else {
						// Single listener in the map. Remove map entry
						deleteList = append(deleteList, eType)
					}
					logger.Debugf("Event Listener - '%s' successfully unregistered for event type - '%s'", name, eType)
					index = -1
				}
			}
		}
	} else {
		for eType, elList := range eventListeners {
			for i, el := range elList {
				if strings.EqualFold(el.Name(), name) {
					index = i
					break
				}
			}
			if index > -1 {
				if len(elList) > 1 {
					// More than one listeners. Just adjust slice
					eventListeners[eType] = append(eventListeners[eType][:index], eventListeners[eType][index+1:]...)
				} else {
					// Single listener in the map. Remove map entry
					deleteList = append(deleteList, eType)
				}
				logger.Debugf("Event Listener - '%s' successfully unregistered for event type - '%s'", name, eType)
				index = -1
			}
		}
	}

	if len(deleteList) > 0 {
		for _, evtType := range deleteList {
			delete(eventListeners, evtType)
		}
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

//  EventContext is a wrapper over specific event
type EventContext struct {
	// Type of event
	eventType string
	// Event data
	event interface{}
}

// Returns wrapped event data
func (ec *EventContext) GetEvent() interface{} {
	return ec.event
}

// Returns event type
func (ec *EventContext) GetType() string {
	return ec.eventType
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
	regListeners, ok := eventListeners[fe.eventType]
	if ok {
		for _, ls := range regListeners {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf("Registered event handler - '%s' failed to process event due to error - '%v' ", ls.Name(), r)
						logger.Errorf("StackTrace: %s", debug.Stack())
					}
				}()
				err := ls.HandleEvent(fe)
				if err != nil {
					logger.Errorf("Registered event handler - '%s' failed to process event due to error - '%s' ", ls.Name(), err.Error())
				} else {
					logger.Debugf("Event - '%s' is successfully delivered to event handler - '%s'", fe.eventType, ls.Name())
				}
			}()
		}
	}
}

//TODO channel to be passed to actions
// Puts event with given type and data on the channel
func PublishEvent(eType string, event interface{}) {

	if len(eventListeners[eType]) > 0 {
		evtContext := &EventContext{event: event, eventType: eType}
		// Put event on the queue
		eventQueue <- evtContext
	}
}
