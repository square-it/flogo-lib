package trigger

import (
	"reflect"
	"strings"
	"sync"
)

var (
	triggersMu sync.Mutex
	triggers   = make(map[string]Trigger)
	reg        = &registry{}
)

type Registry interface {
	TriggerMap() map[string]interface{}
}

type registry struct {
}

func GetRegistry() Registry {
	return reg
}

// Register registers the specified trigger
func Register(trigger Trigger) {
	triggersMu.Lock()
	defer triggersMu.Unlock()

	if trigger == nil {
		panic("trigger.Register: trigger is nil")
	}

	id := trigger.Metadata().ID

	if _, dup := triggers[id]; dup {
		panic("trigger.Register: Register called twice for trigger " + id)
	}

	// copy on write to avoid synchronization on access
	newTriggers := make(map[string]Trigger, len(triggers))

	for k, v := range triggers {
		newTriggers[k] = v
	}

	newTriggers[id] = trigger
	triggers = newTriggers
}

// Triggers gets all the registered triggers
func Triggers() []Trigger {

	var curTriggers = triggers

	list := make([]Trigger, 0, len(curTriggers))

	for _, value := range curTriggers {
		list = append(list, value)
	}

	return list
}

//TriggerTypes returns a map of all the registered Triggers where key is the pkg name of the type
func (r *registry) TriggerMap() map[string]interface{} {
	triggerMap := make(map[string]interface{})

	var curTriggers = triggers

	for _, value := range curTriggers {
		AddTrigger(triggerMap, value)
	}

	return triggerMap
}

func AddTrigger(m map[string]interface{}, value interface{}) {
	t := reflect.TypeOf(value)
	pkgPath := t.Elem().PkgPath()
	pkgPath = strings.TrimLeft(pkgPath, "vendor/src/")
	pkgPath = strings.TrimLeft(pkgPath, "vendor/")
	m[pkgPath] = value
}

// Get gets specified trigger
func Get(id string) Trigger {
	//var curTriggers = triggers
	return triggers[id]
}
