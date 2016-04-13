package test

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// TestActivityContext is a dummy AcitivityContext to assist in testing
type TestActivityContext struct {
	FlowIDVal   string
	FlowNameVal string
	TaskNameVal string
	Attrs       map[string]*data.Attribute
}

// NewTestActivityContext creates a new TestActivityContext
func NewTestActivityContext() *TestActivityContext {
	tc := &TestActivityContext{
		FlowIDVal:   "1",
		FlowNameVal: "Test Flow",
		TaskNameVal: "Test Task",
		Attrs:       make(map[string]*data.Attribute),
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
func (c *TestActivityContext) GetAttrType(attrName string) (attrType string, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Type, true
	}

	return "", false
}

// GetAttrValue implements data.Scope.GetAttrValue
func (c *TestActivityContext) GetAttrValue(attrName string) (value string, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Value, true
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
