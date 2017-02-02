package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockFactory struct {
}

func (f *MockFactory) New(id string) Action2 {
	return nil
}

//TestAddFactoryEmptyRef
func TestAddFactoryEmptyRef(t *testing.T) {

	reg := &registry{}

	// Add factory
	err := reg.AddFactory("", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "registry.RegisterFactory: ref is empty", err.Error())
}

//TestAddFactoryNilFactory
func TestAddFactoryNilFactory(t *testing.T) {

	reg := &registry{}

	// Add factory
	err := reg.AddFactory("github.com/mock", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "registry.RegisterFactory: factory is nil", err.Error())
}

//TestAddFactoryDuplicated
func TestAddFactoryDuplicated(t *testing.T) {

	reg := &registry{}
	f := &MockFactory{}

	// Add factory: this time should pass
	err := reg.AddFactory("github.com/mock", f)
	assert.Nil(t, err)
	// Add factory: this time should fail, duplicated
	err = reg.AddFactory("github.com/mock", f)
	assert.NotNil(t, err)
	assert.Equal(t, "registry.RegisterFactory: already registered factory for ref 'github.com/mock'", err.Error())
}

//TestAddFactoryOk
func TestAddFactoryOk(t *testing.T) {

	reg := &registry{}
	f := &MockFactory{}

	// Add factory
	err := reg.AddFactory("github.com/mock", f)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(reg.factories))
}

//TestAddInstanceEmptyId
func TestAddInstanceEmptyId(t *testing.T) {

	reg := &registry{}

	// Add factory
	err := reg.AddInstance("", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "registry.RegisterInstance: id is empty", err.Error())
}

//TestAddInstanceNilInstance
func TestAddInstanceNilInstance(t *testing.T) {

	reg := &registry{}

	// Add instance
	err := reg.AddInstance("myInstanceId", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "registry.RegisterInstance: instance is nil", err.Error())
}

//TestAddInstanceDuplicated
func TestAddInstanceDuplicated(t *testing.T) {

	reg := &registry{}
	i := &ActionInstance{}

	// Add instance: this time should pass
	err := reg.AddInstance("myinstanceId", i)
	assert.Nil(t, err)
	// Add instance: this time should fail, duplicated
	err = reg.AddInstance("myinstanceId", i)
	assert.NotNil(t, err)
	assert.Equal(t, "registry.RegisterInstance: already registered instance for id 'myinstanceId'", err.Error())
}

//TestAddInstanceOk
func TestAddInstanceOk(t *testing.T) {

	reg := &registry{}
	i := &ActionInstance{}

	// Add instance
	err := reg.AddInstance("myinstanceId", i)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(reg.instances))
}
