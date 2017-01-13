package app

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("app")

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
	// Get Registered triggers
	regTriggers := triggerRegistry.TriggerMap()

	// Get Trigger instances from configuration
	configTriggers := m.App.Triggers

	m.Triggers = make(map[string]*trigger.Trigger, len(configTriggers))

	for _, configTrigger := range configTriggers {
		if configTrigger == nil {
			continue
		}
		regTrigger, ok := regTriggers[configTrigger.Ref]
		if !ok {
			return fmt.Errorf("Trigger '%s' not registered", configTrigger.Ref)
		}

		var newInstance trigger.Trigger
		var network bytes.Buffer

		gob.Register(regTrigger)

		enc := gob.NewEncoder(&network)

		err := enc.Encode(&regTrigger)
		if err != nil {
			return fmt.Errorf("Trigger instance creation '%s'", err.Error())
		}

		dec := gob.NewDecoder(&network)
		err = dec.Decode(&newInstance)
		if err != nil {
			return fmt.Errorf("Trigger instance creation '%s'", err.Error())
		}

		m.Triggers[configTrigger.Id] = &newInstance
	}

	return nil
}
