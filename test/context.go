package test

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
)

// TestActivityContext is a dummy AcitivityContext to assist in testing
type TestActivityContext struct {
	FlowIDVal   string
	FlowNameVal string
	TaskNameVal string
	Attrs       map[string]*data.Attribute

	metadta      *activity.Metadata
	inputs       map[string]*data.Attribute
	outputs      map[string]*data.Attribute
}

// NewTestActivityContext creates a new TestActivityContext
func NewTestActivityContext(metadata *activity.Metadata) *TestActivityContext {
	tc := &TestActivityContext{
		FlowIDVal:   "1",
		FlowNameVal: "Test Flow",
		TaskNameVal: "Test Task",
		Attrs:       make(map[string]*data.Attribute),
		inputs: make(map[string]*data.Attribute, len(metadata.Inputs)),
		outputs: make(map[string]*data.Attribute, len(metadata.Outputs)),
	}

	for _,element := range metadata.Inputs {
		tc.inputs[element.Name] = data.NewAttribute(element.Name, element.Type, nil)
	}
	for _,element := range metadata.Outputs {
		tc.outputs[element.Name] = data.NewAttribute(element.Name, element.Type, nil)
	}

	return tc
}

// FlowInstanceID implements activity.Context.FlowInstanceID
func (c *TestActivityContext) FlowInstanceID() string {
	return c.FlowIDVal
}

// FlowName implements activity.Context.FlowName
func (c *TestActivityContext) FlowName() string {
	return c.FlowNameVal
}

// TaskName implements activity.Context.TaskName
func (c *TestActivityContext) TaskName() string {
	return c.TaskNameVal
}

// GetAttrType implements data.Scope.GetAttrType
func (c *TestActivityContext) GetAttrType(attrName string) (attrType data.Type, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Type, true
	}

	return 0, false
}

// GetAttrValue implements data.Scope.GetAttrValue
func (c *TestActivityContext) GetAttrValue(attrName string) (value string, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Value.(string), true
	}

	return "", false
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *TestActivityContext) SetAttrValue(attrName string, value string) {

	attr, found := c.Attrs[attrName]

	if found {
		attr.Value = value
	}
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *TestActivityContext) SetInput(name string, value interface{}) {

	attr, found := c.inputs[name]

	if found {
		attr.Value = value
	} else {
		//error?
	}
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *TestActivityContext) GetInput(name string) interface{} {

	attr, found := c.inputs[name]

	if found {
		return attr.Value
	}

	return nil
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *TestActivityContext) SetOutput(name string, value interface{}) {

	attr, found := c.outputs[name]

	if found {
		attr.Value = value
	} else {
		//error?
	}
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *TestActivityContext) GetOutput(name string) interface{} {

	attr, found := c.outputs[name]

	if found {
		return attr.Value
	}

	return nil
}
