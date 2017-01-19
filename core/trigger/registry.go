package trigger

import (
	"fmt"
	"sync"
)

var (
	triggersMu sync.Mutex
	triggers   = make(map[string]Trigger)
	reg        = &registry{}
)

type Registry interface {
	RegisterFactory(ref string, f Factory) error
	GetFactories() map[string]Factory
	RegisterInstance(id string, trigger Trigger2) error
}

type registry struct {
	factories map[string]Factory
	instances map[string]Trigger2
}

func GetRegistry() Registry {
	return reg
}

func (r *registry) RegisterFactory(ref string, f Factory) error {
	triggersMu.Lock()
	defer triggersMu.Unlock()

	if len(ref) == 0 {
		return fmt.Errorf("registry.RegisterFactory: ref is empty")
	}

	if f == nil {
		return fmt.Errorf("registry.RegisterFactory: factory is nil")
	}

	// copy on write to avoid synchronization on access
	newFs := make(map[string]Factory, len(r.factories))

	for k, v := range r.factories {
		newFs[k] = v
	}

	if newFs[ref] != nil {
		return fmt.Errorf("registry.RegisterFactory: already registered factory for ref '%s'", ref)
	}

	newFs[ref] = f

	r.factories = newFs

	return nil
}

// GetFactories returns a copy of the factories map
func (r *registry) GetFactories() map[string]Factory {

	newFs := make(map[string]Factory, len(r.factories))

	for k, v := range r.factories {
		newFs[k] = v
	}

	return newFs
}

func (r *registry) RegisterInstance(id string, inst Trigger2) error {
	triggersMu.Lock()
	defer triggersMu.Unlock()

	if len(id) == 0 {
		return fmt.Errorf("registry.RegisterInstance: id is empty")
	}

	if inst == nil {
		return fmt.Errorf("registry.RegisterInstance: instance is nil")
	}

	// copy on write to avoid synchronization on access
	newInst := make(map[string]Trigger2, len(r.instances))

	for k, v := range r.instances {
		newInst[k] = v
	}

	if newInst[id] != nil {
		return fmt.Errorf("registry.RegisterInstance: already registered instance for id '%s'", id)
	}

	newInst[id] = inst

	r.instances = newInst

	return nil

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

// Get gets specified trigger
func Get(id string) Trigger {
	//var curTriggers = triggers
	return triggers[id]
}
