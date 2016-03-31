package test

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type TestActivityContext struct {
	ProcessIDVal string
	ProcessNameVal string
	TaskNameVal string
	Attrs   map[string]*data.Attribute
}

func NewTestActivityContext() *TestActivityContext {
	tc := &TestActivityContext{
		ProcessIDVal:"1",
		ProcessNameVal:"Test Process",
		TaskNameVal:"Test Task",
		Attrs: make(map[string]*data.Attribute),
	}

	return tc
}

// ProcessInstanceID implements activity.Context.ProcessInstanceID
func (c *TestActivityContext) ProcessInstanceID() string {
	return c.ProcessIDVal
}

// ProcessName implements activity.Context.ProcessName
func (c *TestActivityContext) ProcessName() string{
	return c.ProcessNameVal
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


