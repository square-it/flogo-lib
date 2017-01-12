package app

import (
	"fmt"
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/types"
)

//InstanceManager will create and maintain all the trigger and action instances for an app
type InstanceManager struct {
	App      *types.App
	Triggers map[string]*trigger.Trigger
}

//NewInstanceManager creates a new instance manager
func NewInstanceManager(app *types.App) *InstanceManager {
	return &InstanceManager{App: app}
}

//CreateInstances creates new instances for triggers and actions in the registry
func (m *InstanceManager) CreateInstances(triggerRegistry trigger.Registry) error {
	// Get Registered trigger types
	regTriggers := triggerRegistry.TriggerTypes()

	// Get Trigger instances from configuration
	configTriggers := m.App.Triggers

	m.Triggers = make(map[string]*trigger.Trigger, len(configTriggers))

	for _, configTrigger := range configTriggers {
		if configTrigger == nil {
			continue
		}
		regTriggerType, ok := regTriggers[configTrigger.Ref]
		if !ok {
			return fmt.Errorf("Trigger '%s' not registered", configTrigger.Ref)
		}
		instanceValue := reflect.New(regTriggerType)
		// Cast to Trigger
		instance, ok := instanceValue.Elem().Interface().(trigger.Trigger)
		if !ok {
			return fmt.Errorf("Trigger '%s' does not implement trigger.Trigger interface", configTrigger.Ref)
		}
		m.Triggers[configTrigger.Id] = &instance
	}

	return nil
}
