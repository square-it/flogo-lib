package app

import (
	"reflect"
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
func (t *MockTrigger) Init(config *trigger.Config, actionRunner action.Runner) {
	//Noop
}
func (t *MockTrigger) Start() error {
	return nil
}
func (t *MockTrigger) Stop() error {
	return nil
}

//getMockApp returns a mock app
func getMockApp() *types.App {
	triggers := make([]*types.Trigger, 1)

	trigger1 := &types.Trigger{Id: "myTrigger1", Ref: "github.com/TIBCOSoftware/flogo-lib/app"}
	triggers[0] = trigger1

	return &types.App{Name: "MyApp", Version: "1.0.0", Triggers: triggers}
}

type mockTriggerRegistry struct {
}

func (r *mockTriggerRegistry) TriggerTypes() map[string]reflect.Type {
	t := make(map[string]reflect.Type, 1)
	trigger.AddTriggerType(t, &MockTrigger{})
	return t
}

//getMockTriggerRegistry returns a mock trigger registry
func getMockTriggerRegistry() trigger.Registry {
	return &mockTriggerRegistry{}
}
