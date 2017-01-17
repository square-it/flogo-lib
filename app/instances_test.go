package app

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/stretchr/testify/assert"
)

//TestNewEngineErrorNoApp
func TestCreateInstancesOk(t *testing.T) {

	app := getMockApp()
	instanceManager := NewInstanceManager(app)

	err := instanceManager.CreateInstances(getMockTriggerRegistry())

	assert.Nil(t, err)
}

//MockTrigger
type MockTrigger struct {
}

func (t *MockTrigger) Metadata() *trigger.Metadata {
	return nil
}
func (t *MockTrigger) Init(config types.TriggerConfig, actionRunner action.Runner) {
	//Noop
}
func (t *MockTrigger) Start() error {
	return nil
}
func (t *MockTrigger) Stop() error {
	return nil
}
func (t *MockTrigger) New(id string) trigger.Trigger2 {
	return nil
}

//getMockApp returns a mock app
func getMockApp() *types.AppConfig {
	triggers := make([]*types.TriggerConfig, 1)

	trigger1 := &types.TriggerConfig{Id: "myTrigger1", Ref: "github.com/TIBCOSoftware/flogo-lib/app"}
	triggers[0] = trigger1

	return &types.AppConfig{Name: "MyApp", Version: "1.0.0", Triggers: triggers}
}

type mockTriggerRegistry struct {
}

func (r *mockTriggerRegistry) GetTriggers() map[string]trigger.Trigger2 {
	t := make(map[string]trigger.Trigger2, 1)
	trigger.AddTrigger(t, &MockTrigger{})
	return t
}

func (r *mockTriggerRegistry) Add(t trigger.Trigger2) error {
	return nil
}

//getMockTriggerRegistry returns a mock trigger registry
func getMockTriggerRegistry() trigger.Registry {
	return &mockTriggerRegistry{}
}
