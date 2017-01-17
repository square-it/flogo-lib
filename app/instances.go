package app

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("app")

//InstanceManager will create and maintain all the trigger and action instances for an app
type InstanceManager struct {
	App      *types.AppConfig
	Triggers map[string]*TriggerInstance
	Actions  map[string]*ActionInstance
}

//TriggerInstance contains all the information for a Trigger Instance, configuration and interface
type TriggerInstance struct {
	Config    *types.TriggerConfig
	Interface trigger.Trigger2
}

//ActionInstance contains all the information for an Action Instance, configuration and interface
type ActionInstance struct {
	Config    *types.ActionConfig
	Interface action.Action2
}

//NewInstanceManager creates a new instance manager
func NewInstanceManager(app *types.AppConfig) *InstanceManager {
	return &InstanceManager{App: app}
}

//CreateInstances creates new instances for triggers and actions in the registry
func (m *InstanceManager) CreateInstances(triggerRegistry trigger.Registry, actionRegistry action.Registry) error {
	// Create Triggers
	err := m.CreateTriggerInstances(triggerRegistry)
	if err != nil {
		return err
	}

	// Create Actions
	err = m.CreateActionInstances(actionRegistry)
	if err != nil {
		return err
	}

	return nil
}

//CreateTriggerInstances creates new instances for triggers in the registry
func (m *InstanceManager) CreateTriggerInstances(triggerRegistry trigger.Registry) error {
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

		newInterface := regTrigger.New(configTrigger.Id)

		if newInterface == nil {
			return fmt.Errorf("Cannot create Trigger nil for id '%s'", configTrigger.Id)
		}

		m.Triggers[configTrigger.Id] = &TriggerInstance{Config: configTrigger, Interface: newInterface}
	}

	return nil
}

//CreateActionInstances creates new instances for actions in the registry
func (m *InstanceManager) CreateActionInstances(actionRegistry action.Registry) error {
	// Get Registered actions
	regActions := actionRegistry.GetActions()

	// Get Action instances from configuration
	configActions := m.App.Actions

	m.Actions = make(map[string]*ActionInstance, len(configActions))

	for _, configAction := range configActions {
		if configAction == nil {
			continue
		}
		regAction, ok := regActions[configAction.Ref]
		if !ok {
			return fmt.Errorf("Action '%s' not registered", configAction.Ref)
		}

		newInterface := regAction.New(configAction.Id)

		if newInterface == nil {
			return fmt.Errorf("Cannot create Action nil for id '%s'", configAction.Id)
		}

		m.Actions[configAction.Id] = &ActionInstance{Config: configAction, Interface: newInterface}
	}

	return nil
}
