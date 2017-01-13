package trigger

import (
	"fmt"
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
	Add(t Trigger2) error
	GetTriggers() map[string]Trigger2
}

type registry struct {
	triggers map[string]Trigger2
}

func GetRegistry() Registry {
	return reg
}

func (r *registry) Add(t Trigger2) error {
	triggersMu.Lock()
	defer triggersMu.Unlock()

	if t == nil {
		return fmt.Errorf("trigger.Register: trigger is nil")
	}

	if t.Metadata() == nil {
		return fmt.Errorf("trigger.Register: trigger metadata is nil")
	}
	id := t.Metadata().ID

	if _, dup := r.triggers[id]; dup {
		return fmt.Errorf("trigger.Register: Register called twice for trigger '%s' ", id)
	}

	// copy on write to avoid synchronization on access
	newTs := make(map[string]Trigger2, len(r.triggers))

	for k, v := range r.triggers {
		newTs[k] = v
	}

	newTs[id] = t
	r.triggers = newTs

	return nil
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

//GetTriggers returns a map of all the registered Triggers where key is the pkg name of the type
func (r *registry) GetTriggers() map[string]Trigger2 {
	triggerMap := make(map[string]Trigger2)

	var curTriggers = r.triggers

	for _, value := range curTriggers {
		AddTrigger(triggerMap, value)
	}

	return triggerMap
}

func AddTrigger(m map[string]Trigger2, trigger Trigger2) {
	t := reflect.TypeOf(trigger)
	pkgPath := t.Elem().PkgPath()
	pkgPath = strings.TrimLeft(pkgPath, "vendor/src/")
	pkgPath = strings.TrimLeft(pkgPath, "vendor/")
	m[pkgPath] = trigger
}

// Get gets specified trigger
func Get(id string) Trigger {
	//var curTriggers = triggers
	return triggers[id]
}
