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
	App      *types.AppConfig
	Triggers map[string]*TriggerInstance
}

//TriggerInstance contains all the information for a Trigger Instance, configuration and interface
type TriggerInstance struct {
	Config    *types.TriggerConfig
	Interface trigger.Trigger2
}

//NewInstanceManager creates a new instance manager
func NewInstanceManager(app *types.AppConfig) *InstanceManager {
	return &InstanceManager{App: app}
}

//CreateInstances creates new instances for triggers and actions in the registry
func (m *InstanceManager) CreateInstances(triggerRegistry trigger.Registry) error {
	// Get Registered triggers
	regTriggers := triggerRegistry.GetTriggers()

	// Get Trigger instances from configuration
	configTriggers := m.App.Triggers

	m.Triggers = make(map[string]*TriggerInstance, len(configTriggers))

	for _, configTrigger := range configTriggers {
		if configTrigger == nil {
			continue
		}
		regTrigger, ok := regTriggers[configTrigger.Ref]
		if !ok {
			return fmt.Errorf("Trigger '%s' not registered", configTrigger.Ref)
		}

		var newInterface trigger.Trigger2
		var network bytes.Buffer

		gob.Register(regTrigger)

		enc := gob.NewEncoder(&network)

		err := enc.Encode(&regTrigger)
		if err != nil {
			return fmt.Errorf("Trigger instance creation encoding '%s'", err.Error())
		}

		dec := gob.NewDecoder(&network)
		err = dec.Decode(&newInterface)
		if err != nil {
			return fmt.Errorf("Trigger instance creation decoding '%s'", err.Error())
		}

		m.Triggers[configTrigger.Id] = &TriggerInstance{Config: configTrigger, Interface: newInterface}
	}

	return nil
}
