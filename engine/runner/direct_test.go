package runner

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type MockAction struct {
	mock.Mock
}

func (m *MockAction) Config() *action.Config {
	return nil
}

func (m *MockAction) Metadata() *action.Metadata {
	return nil
}

func (m *MockAction) Run(context context.Context, inputs []*data.Attribute, options map[string]interface{}, handler action.ResultHandler) error {
	args := m.Called(context, inputs, options, handler)
	if handler != nil {
		resultData, _ := data.CoerceToObject("{\"data\":\"mock\"}")
		resultData["code"] = 200
		handler.HandleResult(resultData, nil)
		handler.Done()
	}
	return args.Error(0)
}

//Test that Result returns the expected values
func TestResultOk(t *testing.T) {

	resultData, _ := data.CoerceToObject("{\"data\":\"mock data\"}")
	resultData["code"] = 1
	rh := &SyncResultHandler{data: resultData, err: errors.New("New Error")}
	data, err := rh.Result()
	assert.Equal(t, 1, data["code"])
	assert.Equal(t, "mock data", data["data"])
	assert.NotNil(t, err)
}

//Test Direct Start method
func TestDirectStartOk(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	err := runner.Start()
	assert.Nil(t, err)
}

//Test Stop method
func TestDirectStopOk(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	err := runner.Stop()
	assert.Nil(t, err)
}

//Test Run method with a nil action
func TestDirectRunNilAction(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	_, _, err := runner.Run(nil, nil, "", nil)
	assert.NotNil(t, err)
}

//Test Run method with error running action
func TestDirectRunErr(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAction)
	mockAction.On("Run", nil, mock.AnythingOfType("[]*data.Attribute"), mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(errors.New("Action Error"))
	_, _, err := runner.Run(nil, mockAction, "", nil)
	assert.NotNil(t, err)
}

//Test Run method ok
func TestDirectRunOk(t *testing.T) {
	runner := NewDirect()
	assert.NotNil(t, runner)
	// Mock Action
	mockAction := new(MockAction)

	mockAction.On("Run", nil, mock.AnythingOfType("[]*data.Attribute"), mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*runner.SyncResultHandler")).Return(nil)
	code, data, err := runner.Run(nil, mockAction, "", nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, "mock", data)
}
