package app

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("app")

//InstanceManager will create and register all the trigger and action instances for an app
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
func (m *InstanceManager) CreateInstances(tRegistry trigger.Registry, aRegistry action.Registry) error {
	// Create Triggers
	err := m.CreateTriggerInstances(tRegistry)
	if err != nil {
		return err
	}

	// Create Actions
	err = m.CreateActionInstances(aRegistry)
	if err != nil {
		return err
	}

	return nil
}

//CreateTriggerInstances creates new instances for triggers in the registry
func (m *InstanceManager) CreateTriggerInstances(tRegistry trigger.Registry) error {
	// Get Registered trigger factories
	factories := tRegistry.GetFactories()

	// Get Trigger instances from configuration
	triggers := m.App.Triggers

	m.Triggers = make(map[string]*TriggerInstance, len(triggers))

	for _, trigger := range triggers {
		if trigger == nil {
			continue
		}
		factory, ok := factories[trigger.Ref]
		if !ok {
			return fmt.Errorf("Trigger '%s' not registered", trigger.Ref)
		}

		newInterface := factory.New(trigger.Id)

		if newInterface == nil {
			return fmt.Errorf("Cannot create Trigger nil for id '%s'", trigger.Id)
		}

		m.Triggers[trigger.Id] = &TriggerInstance{Config: trigger, Interface: newInterface}
	}

	return nil
}

//CreateActionInstances creates new instances for actions in the registry
func (m *InstanceManager) CreateActionInstances(aRegistry action.Registry) error {
	// Get Registered actions
	regActions := aRegistry.GetActions()

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
