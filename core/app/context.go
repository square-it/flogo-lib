package app

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"sync"
)

type Context interface {
	data.Scope

	//AddAttr add a new data.Attribute to the global application context
	AddAttr(attrName string, attrType data.Type, value interface{})
}

var (
	appCtxt = newAppContext()
)

func GetContext() Context {
	return appCtxt;
}

type appContext struct {
	mutex sync.Mutex
	attrs   map[string]*data.Attribute
}

// NewAppContext creates a new AppContext, there should only be one in an engine
func newAppContext() *appContext {

	context := &appContext{
		attrs:       make(map[string]*data.Attribute),
	}

	return context
}

// GetAttrType implements Scope.GetAttrType
func (ctx *appContext) GetAttrType(attrName string) (attrType data.Type, exists bool) {

	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	attr, found := ctx.attrs[attrName]

	if found {
		return attr.Type, true
	}

	return 0, false
}

// GetAttrValue implements Scope.GetAttrValue
func (ctx *appContext) GetAttrValue(attrName string) (value interface{}, exists bool) {

	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	attr, found := ctx.attrs[attrName]

	if found {
		return attr.Value, true
	}

	return nil, false
}

// SetAttrValue implements Scope.SetAttrValue
func (ctx *appContext) SetAttrValue(attrName string, value interface{}) {

	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	attr, found := ctx.attrs[attrName]

	if found {
		attr.Value = value
	}
	//todo return error? how do we determine type
}

// AddAttr implements Context.AddAttr
func (ctx *appContext) AddAttr(attrName string, attrType data.Type, value interface{}) {

	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	attr, found := ctx.attrs[attrName]

	if found {
		// add type is the same
		attr.Value = value
	}
	//todo return error? how do we determine type
}


